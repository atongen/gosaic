package controller

import (
	"errors"
	"gosaic/environment"
	"gosaic/model"
	"gosaic/service"
	"gosaic/util"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"gopkg.in/cheggaaa/pb.v1"
)

type addIndex struct {
	path   string
	md5sum string
}

func Index(env environment.Environment, paths []string) error {
	found := indexGetPaths(env.Log(), paths)
	if len(found) == 0 {
		return nil
	}

	gidxService, err := env.GidxService()
	if err != nil {
		env.Printf("Error getting index service: %s\n", err.Error())
		return err
	}

	aspectService, err := env.AspectService()
	if err != nil {
		env.Printf("Error getting aspect service: %s\n", err.Error())
		return err
	}

	err = processIndexPaths(env.Log(), env.Workers(), found, gidxService, aspectService)
	if err != nil {
		env.Printf("Error indexing images: %s\n", err.Error())
		return err
	}

	return nil
}

func indexGetPaths(l *log.Logger, paths []string) []string {
	found := make([]string, 0)

	exts := []string{".jpg", ".jpeg", ".png"}
	for _, myPath := range paths {
		err := filepath.Walk(myPath, func(path string, f os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !f.Mode().IsRegular() {
				return nil
			}
			ext := strings.ToLower(filepath.Ext(path))
			if !util.SliceContainsString(exts, ext) {
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

func processIndexPaths(l *log.Logger, workers int, paths []string, gidxService service.GidxService, aspectService service.AspectService) error {
	num := len(paths)
	if num == 0 {
		return nil
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

	cancel := false
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cancel = true
	}()

	for _, p := range paths {
		if cancel {
			break
		}
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
	close(c)
	close(add)
	close(sem)
	close(done)

	if cancel {
		return errors.New("Cancelled")
	}

	bar.Finish()
	return nil
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
		Width:       width,
		Height:      height,
		Orientation: orientation,
	}
	err = gidxService.Insert(&gidx)
	if err != nil {
		return err
	}

	return nil
}
