package command

import (
  "os"
  "github.com/codegangsta/cli"
  "github.com/atongen/gosaic/service"
  "github.com/atongen/gosaic/util"
  "github.com/atongen/gosaic/model"
  "fmt"
  "path/filepath"
  "strings"
)

func IndexAction(env *Environment, c *cli.Context) {
  path := c.Args()[0]
	finfo, err := os.Stat(path)
	if err != nil {
		env.Log.Fatalln("File or directory does not exist: %s\n", path)
	}

	gidxService := service.NewGidxService(env.DB)

	if finfo.IsDir() {
		err = indexDir(gidxService, path)
	} else {
		err = indexFile(gidxService, path)
	}
  if err != nil {
    env.Log.Fatalln(err)
  }
}

func indexDir(gidxService *service.GidxService, path string) error {
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		return indexFile(gidxService, path)
	})
	return err
}

func indexFile(gidxService *service.GidxService, path string) error {
	if strings.HasSuffix(strings.ToLower(path), ".jpg") {
		hash, err := util.Md5sum(path)
		if err != nil {
			return err
		}

		exists, err := gidxService.ExistsByMd5sum(hash)
		if err != nil {
			return err
		}

		if exists {
			return nil
		}

		fmt.Println(path)
		gidx := model.NewGidx(path, hash, 1, 1)
		err = gidxService.Create(gidx)
		if err != nil {
			return err
		}
	}
	return nil
}
