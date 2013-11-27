package runner

import (
	"fmt"
	"github.com/atongen/gosaic/model"
	"github.com/atongen/gosaic/service"
	"github.com/atongen/gosaic/util"
	"os"
	"path/filepath"
	"strings"
)

type Index Run

func (index Index) Execute() error {
	finfo, err := os.Stat(index.Arg)
	if err != nil {
		return fmt.Errorf("File or directory does not exist: %s\n", index.Arg)
	}

	gidxService := service.NewGidxService(index.Project.DB)

	if finfo.IsDir() {
		return indexDir(gidxService, index.Arg)
	} else {
		return indexFile(gidxService, index.Arg)
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
