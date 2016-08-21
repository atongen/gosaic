package controller

import (
	"gosaic/environment"
)

func Status(env environment.Environment) {
	env.Println("Gosaic project db:", env.DbPath())

	gidxService, err := env.GidxService()
	if err != nil {
		env.Println(err.Error())
		return
	}

	count, err := gidxService.Count()
	if err != nil {
		env.Println("Unable to count index")
		return
	}

	env.Println(count, "images in the index.")
	env.Println("Status: OK")
}
