package environment

import (
	"database/sql"
	"io"
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
	DBMEM = ":memory:"
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
	GidxPartialService() (service.GidxPartialService, error)
	CoverService() (service.CoverService, error)
	CoverPartialService() (service.CoverPartialService, error)
	MacroService() (service.MacroService, error)
	MacroPartialService() (service.MacroPartialService, error)
	PartialComparisonService() (service.PartialComparisonService, error)
	MosaicService() (service.MosaicService, error)
	MosaicPartialService() (service.MosaicPartialService, error)
	DbPath() string
	Workers() int
	Log() *log.Logger
	Db() *sql.DB
	DbMap() *gorp.DbMap
	Printf(format string, a ...interface{})
	Println(a ...interface{})
	Fatalf(format string, a ...interface{})
	Fatalln(a ...interface{})
}

type environment struct {
	dbPath   string
	workers  int
	dB       *sql.DB
	dbMap    *gorp.DbMap
	log      *log.Logger
	services map[ServiceName]service.Service
	m        sync.Mutex
}

func NewEnvironment(dbPath string, out io.Writer, workers int) (Environment, error) {
	env := &environment{}

	// setup the environment logger
	env.log = log.New(out, "GOSAIC: ", log.Ldate|log.Ltime)

	var dbPathAbs string
	if dbPath == DBMEM {
		env.dbPath = dbPath
	} else {
		dir := filepath.Dir(dbPath)

		// ensure environment dir exists
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return nil, err
		}

		// get environment absolute dbPath
		dbPathAbs, err = filepath.Abs(dbPath)
		if err != nil {
			return nil, err
		}

		env.dbPath = dbPathAbs
	}

	env.workers = workers

	return env, nil
}

func GetProdEnv(dbPath string, workers int) (Environment, error) {
	return NewEnvironment(dbPath, os.Stdout, workers)
}

func GetTestEnv(out io.Writer) (Environment, error) {
	return NewEnvironment(":memory:", out, 2)
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

func (env *environment) DbPath() string {
	return env.dbPath
}

func (env *environment) Workers() int {
	return env.workers
}

func (env *environment) Log() *log.Logger {
	return env.log
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
