package controller

import (
	"database/sql"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/atongen/gosaic/database"
	_ "github.com/mattn/go-sqlite3"
)

const (
	DbFile = "gosaic.sqlite3"
)

type Environment struct {
	Path    string
	DB      *sql.DB
	DbPath  string
	Log     *log.Logger
	Workers int
	Verbose bool
}

func NewEnvironment(path string, out io.Writer, dbPath string, workers int, verbose bool) *Environment {
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
	env.Verbose = verbose

	return env
}

func GetEnvironment(path string, workers int, verbose bool) *Environment {
	env := NewEnvironment(path, os.Stdout, filepath.Join(path, DbFile), workers, verbose)
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
	// this has been moved to commands
	// how to refactor to bring it back here
	//defer env.DB.Close()

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
}

func (env *Environment) Fatalln(v ...interface{}) {
	env.Log.Fatalln(v)
}

func (env *Environment) Println(v ...interface{}) {
	env.Log.Println(v)
}

func (env *Environment) Verboseln(v ...interface{}) {
	if env.Verbose {
		env.Log.Println(v)
	}
}
