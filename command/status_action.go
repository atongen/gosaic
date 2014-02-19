package command

import (
  "github.com/codegangsta/cli"
  "os"
)

func StatusAction(env *Environment, c *cli.Context) {
  env.Log.Printf("gosiac home: %s\n", env.Path)
	_, err := os.Stat(env.DbPath)
	if err == nil {
		env.Log.Printf("Database exists: %s\n", env.DbPath)
	} else {
		env.Log.Fatalln("Error initializing environment.", err)
	}
}
