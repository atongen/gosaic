package controller

import "gosaic/environment"

func IndexRm(env environment.Environment, paths []string) error {
	gidxService, err := env.GidxService()
	if err != nil {
		return err
	}

	for _, path := range paths {
		gidx, err := gidxService.GetOneBy("path", path)
		if err != nil {
			env.Printf("Error finding image %s: %s\n", path, err.Error())
		}
		if gidx != nil {
			_, err := gidxService.Delete(gidx)
			if err != nil {
				env.Printf("Error removing image %s: %s\n", path, err.Error())
			}
		}
	}

	return nil
}
