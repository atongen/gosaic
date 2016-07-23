package controller

import "gosaic/environment"

func IndexList(env environment.Environment) {
	gidxService, err := env.GidxService()
	if err != nil {
		env.Printf("Error getting index service: %s\n", err.Error())
		return
	}

	batchSize := 1000

	for i := 0; ; i++ {
		gidxs, err := gidxService.FindAll("gidx.path ASC", batchSize, batchSize*i)
		if err != nil {
			env.Printf("Error finding indexes: %s\n", err.Error())
			return
		}
		if len(gidxs) == 0 {
			// we are done
			return
		}

		for _, gidx := range gidxs {
			env.Println(gidx.Path)
		}
	}
}
