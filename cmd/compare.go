package cmd

import (
	"github.com/atongen/gosaic/controller"
	"github.com/spf13/cobra"
)

var (
	compareMacroId int
)

func init() {
	addLocalIntFlag(&compareMacroId, "macro-id", "", 0, "Id of macro for comparison", CompareCmd)
	RootCmd.AddCommand(CompareCmd)
}

var CompareCmd = &cobra.Command{
	Use:    "compare",
	Short:  "Build comparisons for macro against index",
	Long:   "Build comparisons for macro against index",
	Hidden: true,
	Run: func(c *cobra.Command, args []string) {
		if compareMacroId == 0 {
			Env.Fatalln("Macro id is required")
		}

		err := Env.Init()
		if err != nil {
			Env.Fatalf("Unable to initialize environment: %s\n", err.Error())
		}
		defer Env.Close()

		controller.Compare(Env, int64(compareMacroId))
	},
}
