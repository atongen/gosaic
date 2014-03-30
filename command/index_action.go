package command

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/atongen/gosaic/model"
	"github.com/atongen/gosaic/service"
	"github.com/atongen/gosaic/util"
	"github.com/codegangsta/cli"
)

type addIndex struct {
	path   string
	md5sum string
}

func IndexAction(env *Environment, c *cli.Context) {
	if !hasExpectedArgs(c.Args(), 1) {
		env.Log.Fatalln("Path argument is required.")
	}

	path := c.Args()[0]
	paths := getPaths(path, env)
	if len(paths) == 0 {
		env.Log.Println("No images found at path", path)
	}

	processPaths(paths, env)
}

func getPaths(path string, env *Environment) []string {
	f, err := os.Stat(path)
	if err != nil {
		env.Log.Fatalln("File or directory does not exist: %s\n", path)
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
	add := make(chan addIndex)
	go storePaths(add, env)

	sem := make(chan bool, env.Workers)
	for _, path := range paths {
		sem <- true
		go func(path string) {
			defer func() { <-sem }()
			md5sum, err := util.Md5sum(path)
			if err != nil {
				env.Log.Println("Unable to get md5 sum for path", path)
				return
			}
			add <- addIndex{path, md5sum}
		}(path)
	}
	for i := 0; i < cap(sem); i++ {
		sem <- true
	}
}

func storePaths(add chan addIndex, env *Environment) {
	gidxService := service.NewGidxService(env.DB)
	for newIndex := range add {
		storePath(newIndex, gidxService, env)
	}
}

func storePath(newIndex addIndex, gidxService *service.GidxService, env *Environment) {
	exists, err := gidxService.ExistsByMd5sum(newIndex.md5sum)
	if err != nil {
		env.Log.Println("Failed to lookup md5sum", newIndex.md5sum, err)
	}

	if exists {
		return
	}

	bounds, err := util.ImageBounds(newIndex.path)
	if err != nil {
		env.Log.Println("Can't get bounds", newIndex.path, err)
	}

	env.Log.Println(newIndex.path)

	gidx := model.NewGidx(newIndex.path, newIndex.md5sum, uint(bounds.Max.X), uint(bounds.Max.Y))
	err = gidxService.Create(gidx)
	if err != nil {
		env.Log.Println("Error storing image data", newIndex.path, err)
	}
}
