package controller

import "gosaic/environment"

func Status(env environment.Environment) {
	env.Println("Gosaic project directory:", env.Path())

	gidxService, err := env.GidxService()
	if err != nil {
		env.Fatalln(err.Error())
	}

	count, err := gidxService.Count()
	if err != nil {
		env.Fatalln("Unable to count index")
	}

	env.Println(count, "images in the index.")
	env.Println("Status: OK")
}
