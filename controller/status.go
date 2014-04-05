package controller

import "os"

func Status(env *Environment) {
	env.Println("Gosaic project directory:", env.Path)
	if env.DbPath == ":memory:" {
		env.Verboseln("Database in memory")
	} else {
		_, err := os.Stat(env.DbPath)
		if err == nil {
			env.Verboseln("Database exists:", env.DbPath)
		} else {
			env.Fatalln("Error initializing environment.", err)
		}
	}
	env.Println("Status: OK")
}
