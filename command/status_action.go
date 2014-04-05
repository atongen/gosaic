package command

import (
	"os"

	"github.com/codegangsta/cli"
)

func StatusAction(env *Environment, c *cli.Context) {
	if !hasExpectedArgs(c.Args(), 0) {
		env.Fatalln("Unexpected arguments present.")
	}

	env.Println("Gosaic project directory:", env.Path)
	_, err := os.Stat(env.DbPath)
	if err == nil {
		env.Verboseln("Database exists:", env.DbPath)
	} else {
		env.Fatalln("Error initializing environment.", err)
	}
	env.Println("Status: OK")
}
