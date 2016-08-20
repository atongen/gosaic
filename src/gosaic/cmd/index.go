package cmd

import (
	"gosaic/controller"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(IndexCmd)
}

var IndexCmd = &cobra.Command{
	Use:   "index PATHS...",
	Short: "Add path(s) to index",
	Long:  "Add path(s) to index",
	Run: func(c *cobra.Command, args []string) {
		paths := make([]string, 0)
		if len(args) > 0 {
			paths = append(paths, args...)
		}

		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			b, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				Env.Fatalf("Unable to read from stdin")
			}

			inPaths := strings.Split(string(b), "\n")
			for _, inp := range inPaths {
				if inp != "" {
					paths = append(paths, inp)
				}
			}
		}

		if len(paths) == 0 {
			Env.Fatalln("paths to index are required")
		}

		err := Env.Init()
		if err != nil {
			Env.Fatalf("Unable to initialize environment: %s\n", err.Error())
		}
		defer Env.Close()

		controller.Index(Env, paths)
	},
}
