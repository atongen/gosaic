package cmd

import (
	"gosaic/controller"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(ListIndexCmd)
}

var ListIndexCmd = &cobra.Command{
	Use:   "list_index",
	Short: "List index entries",
	Long:  "List index entries",
	Run: func(c *cobra.Command, args []string) {
		err := Env.Init()
		if err != nil {
			Env.Fatalf("Unable to initialize environment: %s\n", err.Error())
		}
		defer Env.Close()

		controller.ListIndex(Env)
	},
}
