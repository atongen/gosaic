package command

import (
  "github.com/codegangsta/cli"
  "runtime"
)

// DirFlag is a global flag
func DirFlag() cli.Flag {
  return cli.StringFlag{"dir, d", ".", "Project directory"}
}

// WorkersFlag is a global flag
func WorkersFlag() cli.Flag {
  return cli.IntFlag{"workers, w", runtime.NumCPU(), "Number of worker processes to use"}
}
