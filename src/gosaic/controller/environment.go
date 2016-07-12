package controller

import (
	"database/sql"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/atongen/gosaic/database"
	"github.com/atongen/gosaic/service"
	"github.com/coopernurse/gorp"
	_ "github.com/mattn/go-sqlite3"
)

const (
	DbFile = "gosaic.sqlite3"
)

type Environment interface {
	Init()
	Close()
	GetService(string) service.Service
	Path() string
	Db() *sql.DB
	DbMap() *gorp.DbMap
	Workers() int
	Fatalln(...interface{})
	Println(...interface{})
	Verboseln(...interface{})
}

type environment struct {
	path     string
	dB       *sql.DB
	dbMap    *gorp.DbMap
	dbPath   string
	log      *log.Logger
	workers  int
	verbose  bool
	debug    bool
	services map[string]service.Service
	m        sync.Mutex
}

func NewEnvironment(path string, out io.Writer, dbPath string, workers int, verbose bool, debug bool) Environment {
	env := &environment{}

	// setup the environment logger
	env.log = log.New(out, "GOSAIC: ", log.Ldate|log.Ltime)

	// get environment absolute path
	path, err := filepath.Abs(path)
	if err != nil {
		env.log.Fatalln("Unable to locate environment absolute path.", err)
	}

	// ensure environment path exists
	err = os.MkdirAll(path, os.ModeDir)
	if err != nil {
		env.log.Fatalln("Unable to create environment path.", err)
	}

	env.path = path
	env.dbPath = dbPath
	env.workers = workers
	env.debug = debug
	if debug {
		env.verbose = true
	} else {
		env.verbose = verbose
	}

	return env
}

func GetProdEnv(path string, workers int, verbose bool, debug bool) Environment {
	return NewEnvironment(path, os.Stdout, filepath.Join(path, DbFile), workers, verbose, debug)
}

func GetTestEnv(out io.Writer) Environment {
	dir, err := ioutil.TempDir("", "GOSAIC")
	if err != nil {
		panic(err)
	}
	return NewEnvironment(dir, out, ":memory:", 2, true, false)
}

func (env *environment) Init() {
	runtime.GOMAXPROCS(env.workers)

	// setup the environment db
	db, err := sql.Open("sqlite3", env.dbPath)
	if err != nil {
		env.Fatalln("Unable to create the db.", err)
	}
	env.dB = db

	// test db connection
	err = env.dB.Ping()
	if err != nil {
		env.Fatalln("Unable to connect to the db.", err)
	}

	// migrate the database
	version, err := database.Migrate(env.dB)
	if err != nil {
		env.Fatalln("Unable to update the db.", err)
	} else {
		env.Verboseln("Database is at version", version)
	}

	// setup orm
	env.dbMap = &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}

	if env.debug {
		env.dbMap.TraceOn("[DB]", env.log)
	}

	// services
	env.services = map[string]service.Service{}
}

func (env *environment) Close() {
	env.dB.Close()
}

func (env *environment) GetService(name string) service.Service {
	env.m.Lock()
	defer env.m.Unlock()

	var s service.Service
	if s, ok := env.services[name]; ok {
		return s
	}

	switch name {
	default:
		env.Fatalln("Service " + name + "not found.")
	case "gidx":
		s = service.NewGidxService(env.dbMap)
	case "aspect":
		s = service.NewAspectService(env.dbMap)
	}
	s.Register()
	env.services[name] = s
	return s
}

func (env *environment) Path() string {
	return env.path
}

func (env *environment) Db() *sql.DB {
	return env.dB
}

func (env *environment) DbMap() *gorp.DbMap {
	return env.dbMap
}

func (env *environment) Workers() int {
	return env.workers
}

func (env *environment) Fatalln(v ...interface{}) {
	env.log.Fatalln(v...)
}

func (env *environment) Println(v ...interface{}) {
	env.log.Println(v...)
}

func (env *environment) Verboseln(v ...interface{}) {
	if env.verbose {
		env.log.Println(v...)
	}
}
