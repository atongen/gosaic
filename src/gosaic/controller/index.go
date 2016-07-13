package controller

import (
	"os"
	"path/filepath"
	"strings"

	"gosaic/environment"
	"gosaic/model"
	"gosaic/util"
)

var (
	total    int
	progress int = 0
)

type addIndex struct {
	path   string
	md5sum string
}

func Index(env environment.Environment, path string) {
	paths := getPaths(path, env)
	total = len(paths)
	if total == 0 {
		env.Println("No images found at path", path)
	} else {
		env.Println("Processing", total, "images")
		processPaths(paths, env)
	}
}

func getPaths(path string, env environment.Environment) []string {
	f, err := os.Stat(path)
	if err != nil {
		env.Fatalln("File or directory does not exist: " + path)
	}

	paths := make([]string, 0)

	if f.IsDir() {
		filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
			if !f.IsDir() {
				if strings.HasSuffix(strings.ToLower(path), ".jpg") {
					absPath, err := filepath.Abs(path)
					if err == nil {
						paths = append(paths, absPath)
					}
				}
			}
			return nil
		})
	} else {
		if strings.HasSuffix(strings.ToLower(path), ".jpg") {
			absPath, err := filepath.Abs(path)
			if err == nil {
				paths = append(paths, absPath)
			}
		}
	}

	return paths
}

func processPaths(paths []string, env environment.Environment) {
	add := make(chan addIndex)
	sem := make(chan bool, env.Workers())

	go storePaths(add, sem, env)

	for _, path := range paths {
		sem <- true
		go func(path string) {
			md5sum, err := util.Md5sum(path)
			if err != nil {
				env.Println("Unable to get md5 sum for path", path)
				return
			}
			add <- addIndex{path, md5sum}
		}(path)
	}

	for i := 0; i < cap(sem); i++ {
		sem <- true
	}
}

func storePaths(add <-chan addIndex, sem <-chan bool, env environment.Environment) {
	for newIndex := range add {
		storePath(newIndex, env)
		<-sem
	}
}

func storePath(newIndex addIndex, env environment.Environment) {
	progress++

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
		env.Println(progress, "of", total, newIndex.path, "(exists)")
		return
	}

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

	env.Println(progress, "of", total, newIndex.path)
}
