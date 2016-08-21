package controller

import (
	"gosaic/environment"
	"gosaic/model"
	"gosaic/service"
	"gosaic/util"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/cheggaaa/pb.v1"
)

type addIndex struct {
	path   string
	md5sum string
}

func Index(env environment.Environment, paths []string) {
	gidxService, err := env.GidxService()
	if err != nil {
		env.Fatalf("Error getting index service: %s\n", err.Error())
	}

	aspectService, err := env.AspectService()
	if err != nil {
		env.Fatalf("Error getting aspect service: %s\n", err.Error())
	}

	found := getJpgPaths(env.Log(), paths)
	processIndexPaths(env.Log(), env.Workers(), found, gidxService, aspectService)
}

func getJpgPaths(l *log.Logger, paths []string) []string {
	found := make([]string, 0)

	for _, myPath := range paths {
		err := filepath.Walk(myPath, func(path string, f os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !f.Mode().IsRegular() {
				return nil
			}
			ext := strings.ToLower(filepath.Ext(path))
			if ext != ".jpg" && ext != ".jpeg" {
				return nil
			}
			absPath, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			if !util.SliceContainsString(found, absPath) {
				found = append(found, absPath)
			}
			return nil
		})

		if err != nil {
			l.Printf("Error building indexing for path %s: %s\n", myPath, err.Error())
		}
	}

	return found
}

func processIndexPaths(l *log.Logger, workers int, paths []string, gidxService service.GidxService, aspectService service.AspectService) {
	num := len(paths)
	if num == 0 {
		return
	}

	l.Printf("Indexing %d images...\n", num)

	bar := pb.StartNew(num)

	add := make(chan addIndex)
	sem := make(chan bool, workers)
	done := make(chan bool)

	go func(myLog *log.Logger, myBar *pb.ProgressBar, addCh <-chan addIndex, semCh <-chan bool, doneCh <-chan bool, myGidxService service.GidxService, myAspectService service.AspectService) {
		for {
			select {
			case newIndex := <-addCh:
				err := storeIndexPath(myLog, newIndex, myGidxService, myAspectService)
				if err != nil {
					l.Printf("Error indexing path %s: %s\n", newIndex.path, err.Error())
				}
				myBar.Increment()
				<-semCh
			case <-doneCh:
				return
			}
		}
	}(l, bar, add, sem, done, gidxService, aspectService)

	for _, p := range paths {
		sem <- true
		go func(myLog *log.Logger, myPath string) {
			md5sum, err := util.Md5sum(myPath)
			if err != nil {
				myLog.Printf("Error getting md5 sum for path %s: %s\n", myPath, err.Error())
				return
			}
			add <- addIndex{myPath, md5sum}
		}(l, p)
	}

	for i := 0; i < cap(sem); i++ {
		sem <- true
	}

	done <- true
	close(add)
	close(sem)
	close(done)

	bar.Finish()
}

func storeIndexPath(l *log.Logger, newIndex addIndex, gidxService service.GidxService, aspectService service.AspectService) error {
	exists, err := gidxService.ExistsBy("md5sum", newIndex.md5sum)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	img, err := util.OpenImage(newIndex.path)
	if err != nil {
		return err
	}

	// don't actually fix orientation here, just determine
	// if x and y need to be swapped
	orientation, err := util.GetOrientation(newIndex.path)
	if err != nil {
		return err
	}

	swap := false
	if 4 < orientation && orientation <= 8 {
		swap = true
	}
	if orientation == 0 {
		orientation = 1
	}

	bounds := (*img).Bounds()

	var width, height int
	if swap {
		width = bounds.Max.Y
		height = bounds.Max.X
	} else {
		width = bounds.Max.X
		height = bounds.Max.Y
	}

	aspect, err := aspectService.FindOrCreate(width, height)
	if err != nil {
		return err
	}

	gidx := model.Gidx{
		AspectId:    aspect.Id,
		Path:        newIndex.path,
		Md5sum:      newIndex.md5sum,
		Width:       uint(width),
		Height:      uint(height),
		Orientation: orientation,
	}
	err = gidxService.Insert(&gidx)
	if err != nil {
		return err
	}

	return nil
}
