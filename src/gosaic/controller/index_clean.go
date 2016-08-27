package controller

import (
	"errors"
	"gosaic/environment"
	"gosaic/model"
	"gosaic/util"
	"os"
	"os/signal"
	"syscall"

	"gopkg.in/cheggaaa/pb.v1"
)

func IndexClean(env environment.Environment) (int, error) {
	num := 0

	gidxService, err := env.GidxService()
	if err != nil {
		return num, err
	}

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

	cancel := false
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cancel = true
	}()

	for i := 0; ; i++ {
		if cancel {
			return num, errors.New("Cancelled")
		}

		gidxs, err := gidxService.FindAll("gidx.id", batchSize, batchSize*i)
		if err != nil {
			return num, err
		}
		if len(gidxs) == 0 {
			// we are done
			return num, nil
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

	if len(toRm) == 0 {
		return num, nil
	}

	for _, gidx := range toRm {
		_, err := gidxService.Delete(gidx)
		if err != nil {
			return num, err
		}
		num++
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
