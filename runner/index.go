package runner

import (
	"fmt"
	"os"
)

type Index Run

func (index Index) Execute() error {
  finfo, err := os.Stat(index.Arg)
  if err != nil {
    return fmt.Errorf("File or directory does not exist: %s\n", index.Arg)
  }
  if finfo.IsDir() {
    return indexDir(index.Arg)
  } else {
    return indexFile(index.Arg)
  }
}

func indexDir(path string) error {
  fmt.Printf("Indexing this directory: %s\n", path)
  return nil
}

func indexFile(path string) error {
  fmt.Printf("Indexing this file: %s\n", path)
  return nil
}
