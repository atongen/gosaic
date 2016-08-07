package cmd

import (
	"gosaic/controller"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(MacroCmd)
}

var MacroCmd = &cobra.Command{
	Use:   "macro PATH COVER_NAME",
	Short: "Add macro",
	Long:  "Add macro",
	Run: func(c *cobra.Command, args []string) {
		if len(args) != 2 {
			Env.Fatalln("Macro path and cover name are required")
		}

		if args[0] == "" {
			Env.Fatalln("Macro path is required")
		}

		if args[1] == "" {
			Env.Fatalln("Cover name is required")
		}

		err := Env.Init()
		if err != nil {
			Env.Fatalf("Unable to initialize environment: %s\n", err.Error())
		}
		defer Env.Close()

		controller.Macro(Env, args[0], args[1])
	},
}
