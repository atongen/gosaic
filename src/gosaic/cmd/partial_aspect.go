package cmd

import (
	"gosaic/controller"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(PartialAspectCmd)
}

var PartialAspectCmd = &cobra.Command{
	Use:   "partial_aspect CxR [CxR CxR...]",
	Short: "Build partial aspects for indexed images",
	Long:  "Build partial aspects for indexed images",
	Run: func(c *cobra.Command, args []string) {
		if len(args) == 0 {
			Env.Fatalln("aspect(s) are required")
		}

		dims := make([]int, len(args)*2)

		for i, aspect := range args {
			aspectSlice := strings.Split(aspect, "x")
			if len(aspectSlice) != 2 {
				Env.Fatalln("Aspect must be of the form: CxR, where C=columns (int) and R=rows (int)")
			}

			colStr := aspectSlice[0]
			rowStr := aspectSlice[1]

			cols, err := strconv.Atoi(colStr)
			if err != nil {
				Env.Fatalln("Aspect columns must be an int")
			}

			rows, err := strconv.Atoi(rowStr)
			if err != nil {
				Env.Fatalln("Aspect rows must be an int")
			}

			dims[i*2] = cols
			dims[i*2+1] = rows
		}

		err := Env.Init()
		if err != nil {
			Env.Fatalf("Unable to initialize environment: %s\n", err.Error())
		}
		defer Env.Close()

		controller.PartialAspect(Env, dims...)
	},
}
