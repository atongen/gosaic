package gosaic

import (
	"database/sql"
	"os"
	"path/filepath"
)

const (
	DbFile = "gosaic.sqlite3"
)

type Project struct {
	Path string
	DB   *sql.DB
}

func NewProject(path string) (*Project, error) {
	err := os.MkdirAll(path, os.ModeDir)
	if err != nil {
		return nil, err
	}
	return &Project{Path: path}, nil
}

func (project *Project) DbPath() string {
	return filepath.Join(project.Path, DbFile)
}
