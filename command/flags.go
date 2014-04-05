package command

import (
	"runtime"

	"github.com/codegangsta/cli"
)

// DirFlag is a global flag
func DirFlag() cli.Flag {
	return cli.StringFlag{"dir, d", ".", "Project directory"}
}

// WorkersFlag is a global flag
func WorkersFlag() cli.Flag {
	return cli.IntFlag{"workers, w", runtime.NumCPU(), "Number of worker processes to use"}
}

// VerboseFlag is a global flag
func VerboseFlag() cli.Flag {
	return cli.BoolFlag{"verbose", "Be verbose"}
}
