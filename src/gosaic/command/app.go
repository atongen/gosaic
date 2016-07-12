package command

import (
	"github.com/codegangsta/cli"
)

func App() *cli.App {
	app := cli.NewApp()
	app.Name = "gosaic"
	app.Usage = "create image mosaics"
	app.Version = "0.0.1"

	// Global flags
	app.Flags = []cli.Flag{
		DirFlag(),
		WorkersFlag(),
		VerboseFlag(),
		DebugFlag(),
	}

	// Commands
	app.Commands = []cli.Command{
		Status(),
		Index(),
	}

	return app
}
