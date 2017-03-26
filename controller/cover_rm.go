package controller

import "github.com/atongen/gosaic/environment"

func CoverRm(env environment.Environment, names []string) {
	coverService := env.ServiceFactory().MustCoverService()

	for _, name := range names {
		cover, err := coverService.GetOneBy("name", name)
		if err != nil {
			env.Printf("Error finding cover %s: %s\n", name, err.Error())
			return
		}
		if cover != nil {
			err := coverService.Delete(cover)
			if err != nil {
				env.Printf("Error removing cover %s\n", name)
			}
		}
	}
}
