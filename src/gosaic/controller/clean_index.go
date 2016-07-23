package controller

import (
	"gosaic/environment"
	"gosaic/model"
	"gosaic/util"
	"os"
)

func CleanIndex(env environment.Environment) {
	gidxService, err := env.GidxService()
	if err != nil {
		env.Printf("Error getting index service: %s\n", err.Error())
		return
	}

	batchSize := 1000

	toRm := []*model.Gidx{}

	for i := 0; ; i++ {
		gidxs, err := gidxService.FindAll("gidx.id", batchSize, batchSize*i)
		if err != nil {
			env.Printf("Error finding indexes: %s\n", err.Error())
			return
		}
		if len(gidxs) == 0 {
			// we are done
			return
		}

		for _, gidx := range gidxs {
			rm, err := shouldRmGidx(gidx)
			if err != nil {
				env.Printf("Error checking index: %s\n", err.Error())
			} else if rm {
				toRm = append(toRm, gidx)
			}
		}
	}

	if len(toRm) == 0 {
		env.Printf("No indexes deleted")
	} else {
		num, err := gidxService.Delete(toRm...)
		if err != nil {
			env.Printf("Error deleting indexes: %s\n", err.Error())
		}
		if num > 0 {
			env.Printf("Deleted %s indexes\n", num)
		}
	}
}

func shouldRmGidx(gidx *model.Gidx) (bool, error) {
	if _, err := os.Stat(gidx.Path); os.IsNotExist(err) {
		return true, nil
	}

	md5sum, err := util.Md5sum(gidx.Path)
	if err != nil {
		return false, err
	}

	if md5sum != gidx.Md5sum {
		return true, nil
	}

	return false, nil
}
