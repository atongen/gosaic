package main

import (
	"github.com/atongen/gosaic/command"
  "github.com/codegangsta/cli"
  "os"
)

func main() {
  app := cli.NewApp()
  app.Name = "gosaic"
  app.Usage = "creates image mosaics"
  app.Version = "0.0.1"

  // Commands
  app.Commands = []cli.Command{
    command.StatusCommand(),
    command.IndexCommand(),
  }

  // Global flags
  app.Flags = []cli.Flag{
    command.DirFlag(),
  }

  app.Run(os.Args)
}
