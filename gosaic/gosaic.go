package main

import (
	"database/sql"
	"flag"
	"github.com/atongen/gosaic"
	"github.com/atongen/gosaic/database"
	"github.com/atongen/gosaic/runner"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

var path string

func init() {
	var defaultPath string
	var err error
	defaultPath, err = os.Getwd()
	if err != nil {
		log.Fatalf("Unable to get current directory.")
	}
	flag.StringVar(&path, "p", defaultPath, "Project path")
}

func main() {
	var run runner.Runner

	flag.Parse()
	subcommand := flag.Arg(0)
	arg := flag.Arg(1)

	// build the project
	project, err := gosaic.NewProject(path)
	if err != nil {
		log.Fatal(err)
	}

	// setup the project db
	db, err := sql.Open("sqlite3", project.DbPath())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	project.DB = db

	// migrate the database
	database.Migrate(project.DB)

	// setup the runner
	switch subcommand {
	case "status":
		run = runner.Status{Project: project, Arg: arg}
	default:
		log.Fatalf("Invalid sub-command: %s\n", subcommand)
	}

	// execute the runner
	err = run.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
