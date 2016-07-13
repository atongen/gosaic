package main

import (
	"os"

	"gosaic/command"
)

func main() {
	command.App().Run(os.Args)
}
