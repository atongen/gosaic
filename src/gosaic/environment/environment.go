package environment

import (
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"

	"gosaic/database"
	"gosaic/service"

	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/gorp.v1"
)

const (
	DbFile = "gosaic.sqlite3"
)

type ServiceName uint8

const (
	GidxServiceName ServiceName = iota
	AspectServiceName
)

var (
	Version   = "none"
	BuildTime = "none"
	BuildUser = "none"
	BuildHash = "none"
)

type Environment interface {
	Init() error
	Close()
	GidxService() (service.GidxService, error)
	AspectService() (service.AspectService, error)
	Path() string
	Workers() int
	Db() *sql.DB
	DbMap() *gorp.DbMap
	Printf(format string, a ...interface{})
	Println(a ...interface{})
	Fatalf(format string, a ...interface{})
	Fatalln(a ...interface{})
}

type environment struct {
	path     string
	workers  int
	dB       *sql.DB
	dbMap    *gorp.DbMap
	dbPath   string
	log      *log.Logger
	services map[ServiceName]service.Service
	m        sync.Mutex
}

func NewEnvironment(path string, out io.Writer, dbPath string, workers int) (Environment, error) {
	env := &environment{}

	// setup the environment logger
	env.log = log.New(out, "GOSAIC: ", log.Ldate|log.Ltime)

	// get environment absolute path
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	// ensure environment path exists
	err = os.MkdirAll(path, os.ModeDir)
	if err != nil {
		return nil, err
	}

	env.path = path
	env.workers = workers
	env.dbPath = dbPath

	return env, nil
}

func GetProdEnv(path string, workers int) (Environment, error) {
	return NewEnvironment(path, os.Stdout, filepath.Join(path, DbFile), workers)
}

func GetTestEnv(out io.Writer) (Environment, error) {
	path, err := ioutil.TempDir("", "GOSAIC")
	if err != nil {
		return nil, err
	}
	return NewEnvironment(path, out, ":memory:", 2)
}

func (env *environment) Init() error {
	// setup the environment db
	db, err := sql.Open("sqlite3", env.dbPath)
	if err != nil {
		return err
	}
	env.dB = db

	// test db connection
	err = env.dB.Ping()
	if err != nil {
		return err
	}

	_, err = env.dB.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return err
	}

	// migrate the database
	_, err = database.Migrate(env.dB)
	if err != nil {
		return err
	}

	// setup orm
	env.dbMap = &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}

	// development
	//env.dbMap.TraceOn("[DB]", env.log)

	// services
	env.services = map[ServiceName]service.Service{}

	return nil
}

func (env *environment) Close() {
	env.dB.Close()
}

func (env *environment) getService(name ServiceName) (service.Service, error) {
	env.m.Lock()
	defer env.m.Unlock()

	var s service.Service
	if s, ok := env.services[name]; ok {
		return s, nil
	}

	switch name {
	default:
		return nil, fmt.Errorf("Service not found")
	case GidxServiceName:
		s = service.NewGidxService(env.dbMap)
	case AspectServiceName:
		s = service.NewAspectService(env.dbMap)
	}
	err := s.Register()
	if err != nil {
		return nil, err
	}
	env.services[name] = s
	return s, nil
}

func (env *environment) GidxService() (service.GidxService, error) {
	s, err := env.getService(GidxServiceName)
	if err != nil {
		return nil, err
	}

	gidxService, ok := s.(service.GidxService)
	if !ok {
		return nil, fmt.Errorf("Invalid gidx service")
	}

	return gidxService, nil
}

func (env *environment) AspectService() (service.AspectService, error) {
	s, err := env.getService(AspectServiceName)
	if err != nil {
		return nil, err
	}

	aspectService, ok := s.(service.AspectService)
	if !ok {
		return nil, fmt.Errorf("Invalid aspect service")
	}

	return aspectService, nil
}

func (env *environment) Path() string {
	return env.path
}

func (env *environment) Workers() int {
	return env.workers
}

func (env *environment) Db() *sql.DB {
	return env.dB
}

func (env *environment) DbMap() *gorp.DbMap {
	return env.dbMap
}

func (env *environment) Fatalln(v ...interface{}) {
	env.log.Fatalln(v...)
}

func (env *environment) Fatalf(format string, v ...interface{}) {
	env.log.Fatalf(format, v...)
}

func (env *environment) Println(v ...interface{}) {
	env.log.Println(v...)
}

func (env *environment) Printf(format string, v ...interface{}) {
	env.log.Printf(format, v...)
}
