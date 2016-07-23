package controller

import "gosaic/environment"

func IndexRm(env environment.Environment, paths []string) {
	gidxService, err := env.GidxService()
	if err != nil {
		env.Printf("Error getting index service: %s\n", err.Error())
		return
	}

	for _, path := range paths {
		gidx, err := gidxService.GetOneBy("path", path)
		if err != nil {
			env.Printf("Error finding indexes: %s\n", err.Error())
			return
		}
		if gidx != nil {
			_, err := gidxService.Delete(gidx)
			if err != nil {
				env.Printf("Error removing %s\n", path)
			}
		}
	}
}
