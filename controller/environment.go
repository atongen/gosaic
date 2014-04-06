package controller

import (
	"database/sql"
	"io"
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

type Environment struct {
	Path    string
	DB      *sql.DB
	DbMap   *gorp.DbMap
	DbPath  string
	Log     *log.Logger
	Workers int
	Verbose bool
	Debug   bool

	services map[string]service.Service
	m        sync.Mutex
}

func NewEnvironment(path string, out io.Writer, dbPath string, workers int, verbose bool, debug bool) *Environment {
	env := &Environment{}

	// setup the environment logger
	env.Log = log.New(out, "GOSAIC: ", log.Ldate|log.Ltime)

	// get environment absolute path
	path, err := filepath.Abs(path)
	if err != nil {
		env.Log.Fatalln("Unable to locate environment absolute path.", err)
	}

	// ensure environment path exists
	err = os.MkdirAll(path, os.ModeDir)
	if err != nil {
		env.Log.Fatalln("Unable to create environment path.", err)
	}

	env.Path = path
	env.DbPath = dbPath
	env.Workers = workers
	env.Debug = debug
	if debug {
		env.Verbose = true
	} else {
		env.Verbose = verbose
	}

	return env
}

func GetEnvironment(path string, workers int, verbose bool, debug bool) *Environment {
	env := NewEnvironment(path, os.Stdout, filepath.Join(path, DbFile), workers, verbose, debug)
	env.Init()
	return env
}

func (env *Environment) Init() {
	runtime.GOMAXPROCS(env.Workers)

	// setup the environment db
	db, err := sql.Open("sqlite3", env.DbPath)
	if err != nil {
		env.Fatalln("Unable to create the db.", err)
	}
	env.DB = db

	// test db connection
	err = env.DB.Ping()
	if err != nil {
		env.Fatalln("Unable to connect to the db.", err)
	}

	// migrate the database
	version, err := database.Migrate(env.DB)
	if err != nil {
		env.Fatalln("Unable to update the db.", err)
	} else {
		env.Verboseln("Database is at version", version)
	}

	// setup orm
	env.DbMap = &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}

	if env.Debug {
		env.DbMap.TraceOn("[DB]", env.Log)
	}

	// services
	env.services = map[string]service.Service{}
}

func (env *Environment) GetService(name string) service.Service {
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
		s = service.NewGidxService(env.DbMap)
		s.Register()
		env.services["gidx"] = s
	}

	return s
}

func (env *Environment) Fatalln(v ...interface{}) {
	env.Log.Fatalln(v...)
}

func (env *Environment) Println(v ...interface{}) {
	env.Log.Println(v...)
}

func (env *Environment) Verboseln(v ...interface{}) {
	if env.Verbose {
		env.Log.Println(v...)
	}
}
