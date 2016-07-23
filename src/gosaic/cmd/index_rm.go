package cmd

import (
	"gosaic/controller"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(IndexRmCmd)
}

var IndexRmCmd = &cobra.Command{
	Use:   "index_rm PATHS...",
	Short: "Remove index entries",
	Long:  "Remove index entries",
	Run: func(c *cobra.Command, args []string) {
		err := Env.Init()
		if err != nil {
			Env.Fatalf("Unable to initialize environment: %s\n", err.Error())
		}
		defer Env.Close()

		controller.IndexRm(Env, args)
	},
}
