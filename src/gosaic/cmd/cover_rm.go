package cmd

import (
	"gosaic/controller"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(CoverRmCmd)
}

var CoverRmCmd = &cobra.Command{
	Use:   "cover_rm NAMES...",
	Short: "Remove cover entries",
	Long:  "Remove cover entries",
	Run: func(c *cobra.Command, args []string) {
		err := Env.Init()
		if err != nil {
			Env.Fatalf("Unable to initialize environment: %s\n", err.Error())
		}
		defer Env.Close()

		controller.CoverRm(Env, args)
	},
}
