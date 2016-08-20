package cmd

import (
	"gosaic/controller"

	"github.com/spf13/cobra"
)

var (
	coverId int
)

func init() {
	addLocalIntFlag(&coverId, "cover_id", "", 0, "Id of cover to use for macro", MacroCmd)
	RootCmd.AddCommand(MacroCmd)
}

var MacroCmd = &cobra.Command{
	Use:   "macro PATH",
	Short: "Add macro",
	Long:  "Add macro",
	Run: func(c *cobra.Command, args []string) {
		if len(args) != 1 {
			Env.Fatalln("Macro path is required")
		}

		if args[0] == "" {
			Env.Fatalln("Macro path is required")
		}

		if coverId == 0 {
			Env.Fatalln("Cover id is required")
		}

		err := Env.Init()
		if err != nil {
			Env.Fatalf("Unable to initialize environment: %s\n", err.Error())
		}
		defer Env.Close()

		controller.Macro(Env, args[0], int64(coverId))
	},
}
