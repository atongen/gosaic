package cmd

import (
	"github.com/atongen/gosaic/controller"
	"github.com/spf13/cobra"
)

var (
	partialAspectMacroId    int
	partialAspectThreashold float64
)

func init() {
	addLocalIntFlag(&partialAspectMacroId, "macro-id", "", 0, "Id of macro to build partials", PartialAspectCmd)
	addLocalFloatFlag(&partialAspectThreashold, "threashold", "t", -1.0, "How similar aspect ratios must be", PartialAspectCmd)
	RootCmd.AddCommand(PartialAspectCmd)
}

var PartialAspectCmd = &cobra.Command{
	Use:    "partial_aspect",
	Short:  "Build partial aspects for indexed images",
	Long:   "Build partial aspects for indexed images",
	Hidden: true,
	Run: func(c *cobra.Command, args []string) {
		err := Env.Init()
		if err != nil {
			Env.Fatalf("Unable to initialize environment: %s\n", err.Error())
		}
		defer Env.Close()

		if partialAspectMacroId == 0 {
			Env.Fatalf("Macro id is required")
		}

		controller.PartialAspect(Env, int64(partialAspectMacroId), partialAspectThreashold)
	},
}
