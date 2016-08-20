package cmd

import (
	"gosaic/controller"

	"github.com/spf13/cobra"
)

func init() {
	addLocalIntFlag(&macroId, "macro_id", "", 0, "Id of macro to build partials", CompareCmd)
	RootCmd.AddCommand(PartialAspectCmd)
}

var PartialAspectCmd = &cobra.Command{
	Use:   "partial_aspect",
	Short: "Build partial aspects for indexed images",
	Long:  "Build partial aspects for indexed images",
	Run: func(c *cobra.Command, args []string) {
		err := Env.Init()
		if err != nil {
			Env.Fatalf("Unable to initialize environment: %s\n", err.Error())
		}
		defer Env.Close()

		if macroId == 0 {
			Env.Fatalf("Macro id is required")
		}

		controller.PartialAspect(Env, int64(macroId))
	},
}
