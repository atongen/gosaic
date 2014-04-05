package main

import (
	"os"

	"github.com/atongen/gosaic/command"
)

func main() {
	command.App().Run(os.Args)
}
