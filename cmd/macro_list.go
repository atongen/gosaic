package cmd

import (
	"github.com/atongen/gosaic/controller"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(MacroListCmd)
}

var MacroListCmd = &cobra.Command{
	Use:    "macro_list",
	Short:  "List macro entries",
	Long:   "List macro entries",
	Hidden: true,
	Run: func(c *cobra.Command, args []string) {
		err := Env.Init()
		if err != nil {
			Env.Fatalf("Unable to initialize environment: %s\n", err.Error())
		}
		defer Env.Close()

		controller.MacroList(Env)
	},
}
