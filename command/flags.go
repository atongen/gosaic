package command

import (
  "github.com/codegangsta/cli"
  "os"
)

func DirFlag() cli.Flag {
  currentDir, err := os.Getwd()
  if err != nil {
    os.Exit(1)
  }
  return cli.StringFlag{"dir, d", currentDir, "directory to place working files"}
}
