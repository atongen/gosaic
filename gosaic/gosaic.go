package main

import (
	"os"

	"github.com/atongen/gosaic/command"
	"github.com/codegangsta/cli"
)

// http://www.easyrgb.com/index.php
// http://en.wikipedia.org/wiki/Color_difference
// http://en.wikipedia.org/wiki/Dithering
// http://en.wikipedia.org/wiki/Color_quantization
func main() {
	app := cli.NewApp()
	app.Name = "gosaic"
	app.Usage = "create image mosaics"
	app.Version = GOSAIC_VERSION

	// Global flags
	app.Flags = []cli.Flag{
		command.DirFlag(),
		command.WorkersFlag(),
		command.VerboseFlag(),
	}

	// Commands
	app.Commands = []cli.Command{
		command.StatusCommand(),
		command.IndexCommand(),
	}

	app.Run(os.Args)
}
