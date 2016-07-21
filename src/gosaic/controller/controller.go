package controller

import (
	"gosaic/environment"
	"os"
	"path/filepath"
	"strings"
)

func getJpgPaths(path string, env environment.Environment) ([]string, error) {
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
