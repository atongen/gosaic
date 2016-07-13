package controller

import (
	"gosaic/service"
)

func Status(env Environment) {
	env.Println("Gosaic project directory:", env.Path())

	gidxService := env.GetService("gidx").(service.GidxService)
	count, err := gidxService.Count()
	if err != nil {
		env.Println("Unable to count index")
	}
	env.Println(count, "images in the index.")
	env.Println("Status: OK")
}
