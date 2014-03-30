package main

import (
	"os"

	"github.com/atongen/gosaic/command"
	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "gosaic"
	app.Usage = "create image mosaics"
	app.Version = GOSAIC_VERSION

	// Global flags
	app.Flags = []cli.Flag{
		command.DirFlag(),
		command.WorkersFlag(),
	}

	// Commands
	app.Commands = []cli.Command{
		command.StatusCommand(),
		command.IndexCommand(),
	}

	app.Run(os.Args)
}
