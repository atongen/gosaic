package cmd

import (
	"gosaic/controller"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(IndexListCmd)
}

var IndexListCmd = &cobra.Command{
	Use:   "index_list",
	Short: "List index entries",
	Long:  "List index entries",
	Run: func(c *cobra.Command, args []string) {
		err := Env.Init()
		if err != nil {
			Env.Fatalf("Unable to initialize environment: %s\n", err.Error())
		}
		defer Env.Close()

		controller.IndexList(Env)
	},
}
