package cmd

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/atongen/gosaic/controller"
	"github.com/spf13/cobra"
)

var (
	indexClean bool
	indexList  bool
	indexRm    bool
)

func init() {
	addLocalBoolFlag(&indexClean, "clean", "c", false, "Clean the index", IndexCmd)
	addLocalBoolFlag(&indexList, "list", "l", false, "List the index", IndexCmd)
	addLocalBoolFlag(&indexRm, "rm", "r", false, "Remove entries from the index", IndexCmd)
	RootCmd.AddCommand(IndexCmd)
}

var IndexCmd = &cobra.Command{
	Use:   "index [PATHS...]",
	Short: "Manage index images",
	Long:  "Manage index images",
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

		err := Env.Init()
		if err != nil {
			Env.Fatalf("Unable to initialize environment: %s\n", err.Error())
		}
		defer Env.Close()

		if indexClean {
			// clean index
			if len(paths) != 0 {
				Env.Fatalln("Cannot specify paths with index clean")
			}
			_, err = controller.IndexClean(Env)
			if err != nil {
				Env.Printf("Error cleaning index: %s\n", err.Error())
			}
		} else if indexList {
			// list index
			if len(paths) != 0 {
				Env.Fatalln("Cannot specify paths with index list")
			}
			err = controller.IndexList(Env)
			if err != nil {
				Env.Printf("Error listing index: %s\n", err.Error())
			}
		} else if indexRm {
			// rm index
			if len(paths) == 0 {
				Env.Fatalln("Must specify paths to rm from index")
			}
			err = controller.IndexRm(Env, paths)
			if err != nil {
				Env.Printf("Error removing index images: %s\n", err.Error())
			}
		} else {
			// add index
			if len(paths) == 0 {
				Env.Fatalln("Must specify paths to index")
			}
			err = controller.Index(Env, paths)
			if err != nil {
				Env.Printf("Error adding index images: %s\n", err.Error())
			}
		}
	},
}
