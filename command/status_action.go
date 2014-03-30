package command

import (
	"os"

	"github.com/codegangsta/cli"
)

func StatusAction(env *Environment, c *cli.Context) {
	if !hasExpectedArgs(c.Args(), 0) {
		env.Log.Fatalln("Unexpected arguments present.")
	}

	env.Log.Printf("Gosaic home: %s\n", env.Path)
	_, err := os.Stat(env.DbPath)
	if err == nil {
		env.Log.Printf("Database exists: %s\n", env.DbPath)
	} else {
		env.Log.Fatalln("Error initializing environment.", err)
	}
}
