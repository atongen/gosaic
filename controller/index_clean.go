package controller

import (
	"errors"
	"github.com/atongen/gosaic/environment"
	"github.com/atongen/gosaic/model"
	"github.com/atongen/gosaic/util"
	"os"

	"gopkg.in/cheggaaa/pb.v1"
)

func IndexClean(env environment.Environment) (int, error) {
	gidxService := env.ServiceFactory().MustGidxService()

	num := 0

	count, err := gidxService.Count()
	if err != nil {
		return num, err
	}

	if count == 0 {
		return num, nil
	}

	env.Printf("Checking %d images...\n", count)

	bar := pb.StartNew(int(count))

	batchSize := 1000
	toRm := []*model.Gidx{}

	for i := 0; ; i++ {
		if env.Cancel() {
			return num, errors.New("Cancelled")
		}

		gidxs, err := gidxService.FindAll("gidx.id", batchSize, batchSize*i)
		if err != nil {
			return num, err
		}
		if len(gidxs) == 0 {
			// we are done
			break
		}

		for _, gidx := range gidxs {
			rm, err := shouldRmGidx(gidx)
			if err != nil {
				return num, err
			} else if rm {
				toRm = append(toRm, gidx)
			}
			bar.Increment()
		}
	}

	if len(toRm) > 0 {
		for _, gidx := range toRm {
			_, err := gidxService.Delete(gidx)
			if err != nil {
				return num, err
			}
			num++
		}
	}

	bar.Finish()
	return num, nil
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
