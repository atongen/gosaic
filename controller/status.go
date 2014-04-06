package controller

import (
	"os"

	"github.com/atongen/gosaic/service"
)

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

	gidxService := env.GetService("gidx").(service.GidxService)
	count, err := gidxService.Count()
	if err != nil {
		env.Println("Unable to count index")
	}
	env.Println(count, "images in the index.")
	env.Println("Status: OK")
}
