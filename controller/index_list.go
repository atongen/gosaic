package controller

import "github.com/atongen/gosaic/environment"

func IndexList(env environment.Environment) error {
	gidxService := env.ServiceFactory().MustGidxService()

	batchSize := 1000

	for i := 0; ; i++ {
		gidxs, err := gidxService.FindAll("gidx.path ASC", batchSize, batchSize*i)
		if err != nil {
			return err
		}
		if len(gidxs) == 0 {
			// we are done
			return nil
		}

		for _, gidx := range gidxs {
			env.Println(gidx.Path)
		}
	}

	return nil
}
