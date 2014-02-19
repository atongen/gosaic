package command

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"os"
  "log"
  "io"
	"path/filepath"
	"github.com/atongen/gosaic/database"
  "runtime"
)

const (
	DbFile = "gosaic.sqlite3"
)

type Environment struct {
	Path        string
	DB          *sql.DB
  DbPath      string
  Log         *log.Logger
  Concurrency int
}

func NewEnvironment(path string, out io.Writer, dbPath string, concurrency int) *Environment {
  env := &Environment{}

  // setup the environment logger
  env.Log = log.New(out, "GOSAIC: ", log.Ldate|log.Ltime|log.Lshortfile)

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
  env.Concurrency = concurrency

	return env
}

func GetEnvironment(path string) *Environment {
  env := NewEnvironment(path, os.Stdout, filepath.Join(DbFile), runtime.NumCPU())
  env.Init()
  return env
}

func (env *Environment) Init() {
  // set concurrency level
  runtime.GOMAXPROCS(env.Concurrency)

	// setup the environment db
	db, err := sql.Open("sqlite3", env.DbPath)
	if err != nil {
		env.Log.Fatalln("Unable to create the db.", err)
	}
	env.DB = db
  defer env.DB.Close()

  // test db connection
  err = env.DB.Ping()
	if err != nil {
		env.Log.Fatalln("Unable to connect to the db.", err)
	}

	// migrate the database
	database.Migrate(env.DB)
}
