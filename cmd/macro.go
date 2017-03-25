package cmd

import (
	"github.com/atongen/gosaic/controller"
	"github.com/spf13/cobra"
)

var (
	macroCoverId int
	macroOutfile string
)

func init() {
	addLocalIntFlag(&macroCoverId, "cover-id", "c", 0, "Id of cover to use for macro", MacroCmd)
	addLocalStrFlag(&macroOutfile, "out", "o", "", "Outfile for resized macro image", MacroCmd)
	RootCmd.AddCommand(MacroCmd)
}

var MacroCmd = &cobra.Command{
	Use:    "macro PATH",
	Short:  "Add macro",
	Long:   "Add macro",
	Hidden: true,
	Run: func(c *cobra.Command, args []string) {
		if len(args) != 1 {
			Env.Fatalln("Macro path is required")
		}

		if args[0] == "" {
			Env.Fatalln("Macro path is required")
		}

		if macroCoverId == 0 {
			Env.Fatalln("Cover id is required")
		}

		err := Env.Init()
		if err != nil {
			Env.Fatalf("Unable to initialize environment: %s\n", err.Error())
		}
		defer Env.Close()

		controller.Macro(Env, args[0], int64(macroCoverId), macroOutfile)
	},
}
