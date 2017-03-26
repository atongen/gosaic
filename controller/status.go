package controller

import (
	"github.com/atongen/gosaic/environment"
)

func Status(env environment.Environment) {
	gidxService := env.ServiceFactory().MustGidxService()

	count, err := gidxService.Count()
	if err != nil {
		env.Println("Unable to count index")
		return
	}

	env.Println(count, "images in the index.")
	env.Println("Status: OK")
}
