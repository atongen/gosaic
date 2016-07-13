package controller

import (
	"os"
	"path/filepath"
	"strings"

	"gosaic/environment"
	"gosaic/model"
	"gosaic/util"
)

type addIndex struct {
	path   string
	md5sum string
}

func Index(env environment.Environment, path string) {
	paths, err := getPaths(path, env)
	if err != nil {
		env.Printf("Error finding images in path %s: %s\n", path, err.Error())
		return
	}

	num := len(paths)
	if num == 0 {
		env.Printf("No images found at path %s\n", path)
		return
	}

	env.Printf("Processing %d images\n", num)
	err = processPaths(paths, env)
	if err != nil {
		env.Printf("Error indexing images: %s\n", err.Error())
	}
}

func getPaths(path string, env environment.Environment) ([]string, error) {
	paths := make([]string, 0)

	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !f.Mode().IsRegular() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".jpg" {
			return nil
		}
		absPath, err := filepath.Abs(path)
		if err != nil {
			return err
		}
		paths = append(paths, absPath)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return paths, nil
}

func processPaths(paths []string, env environment.Environment) error {
	add := make(chan addIndex)
	sem := make(chan bool, env.Workers())

	go storePaths(add, sem, env)

	for _, p := range paths {
		sem <- true
		go func(myPath string) {
			md5sum, err := util.Md5sum(myPath)
			if err != nil {
				env.Printf("Unable to get md5 sum for path %s\n", myPath)
				return
			}
			add <- addIndex{myPath, md5sum}
		}(p)
	}

	for i := 0; i < cap(sem); i++ {
		sem <- true
	}

	return nil
}

func storePaths(add <-chan addIndex, sem <-chan bool, env environment.Environment) {
	for newIndex := range add {
		storePath(newIndex, env)
		<-sem
	}
}

func storePath(newIndex addIndex, env environment.Environment) {
	gidxService, err := env.GidxService()
	if err != nil {
		env.Println(err.Error())
		return
	}

	aspectService, err := env.AspectService()
	if err != nil {
		env.Println(err.Error())
		return
	}

	exists, err := gidxService.ExistsBy("md5sum", newIndex.md5sum)
	if err != nil {
		env.Println("Failed to lookup md5sum", newIndex.md5sum, err)
		return
	}

	if exists {
		return
	}

	env.Println(newIndex.path)

	img, err := util.OpenImage(newIndex.path)
	if err != nil {
		env.Println("Can't open image", newIndex.path, err)
		return
	}

	orientation, err := util.GetOrientation(newIndex.path)
	if err == nil {
		util.FixOrientation(img, orientation)
	}

	bounds := (*img).Bounds()

	aspect, err := aspectService.FindOrCreate(bounds.Max.X, bounds.Max.Y)
	if err != nil {
		env.Println("Error getting image aspect data", newIndex.path, err)
		return
	}

	gidx := model.NewGidx(aspect.Id, newIndex.path, newIndex.md5sum, uint(bounds.Max.X), uint(bounds.Max.Y), orientation)
	err = gidxService.Insert(gidx)
	if err != nil {
		env.Println("Error storing image data", newIndex.path, err)
	}
}
