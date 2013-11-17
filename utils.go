package gosaic

import (
	"os"
	"path/filepath"
)

const (
	DbFile = "gosaic.sqlite3"
)

func DbPath(projectDir string) string {
	return filepath.Join(projectDir, DbFile)
}

func IsProject(projectDir string) (bool, error) {
	dbPath := DbPath(projectDir)
	_, err := os.Stat(dbPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func MkProjectDir(projectDir string) error {
	return os.MkdirAll(projectDir, os.ModeDir)
}
