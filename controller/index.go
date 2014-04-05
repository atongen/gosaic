package controller

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/atongen/gosaic/model"
	"github.com/atongen/gosaic/service"
	"github.com/atongen/gosaic/util"
)

var (
	total    int
	progress int = 0
)

type addIndex struct {
	path   string
	md5sum string
}

func Index(env *Environment, path string) {
	paths := getPaths(path, env)
	total = len(paths)
	if total == 0 {
		env.Println("No images found at path", path)
	} else {
		env.Println("Processing", total, "images")
		processPaths(paths, env)
	}
}

func getPaths(path string, env *Environment) []string {
	f, err := os.Stat(path)
	if err != nil {
		env.Fatalln("File or directory does not exist: %s\n", path)
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

func processPaths(paths []string, env *Environment) {
	gidxService := service.NewGidxService(env.DB)
	add := make(chan addIndex)
	sem := make(chan bool, env.Workers)

	go storePaths(gidxService, add, sem, env)

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

func storePaths(gidxService *service.GidxService, add <-chan addIndex, sem <-chan bool, env *Environment) {
	for newIndex := range add {
		storePath(gidxService, newIndex, env)
		<-sem
	}
}

func storePath(gidxService *service.GidxService, newIndex addIndex, env *Environment) {
	progress++

	exists, err := gidxService.ExistsByMd5sum(newIndex.md5sum)
	if err != nil {
		env.Println("Failed to lookup md5sum", newIndex.md5sum, err)
		return
	}

	if exists {
		env.Verboseln(progress, "of", total, newIndex.path, "(exists)")
		return
	}

	img, err := util.OpenImage(newIndex.path)
	if err != nil {
		env.Println("Can't open image", newIndex.path, err)
		return
	}
	bounds := (*img).Bounds()

	gidx := model.NewGidx(newIndex.path, newIndex.md5sum, uint(bounds.Max.X), uint(bounds.Max.Y))
	err = gidxService.Create(gidx)
	if err != nil {
		env.Println("Error storing image data", newIndex.path, err)
	}

	env.Verboseln(progress, "of", total, newIndex.path)
}
