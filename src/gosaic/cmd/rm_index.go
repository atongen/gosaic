package cmd

import (
	"gosaic/controller"

	"github.com/spf13/cobra"
)

var (
	rmPath string
)

func init() {
	RootCmd.AddCommand(RmIndexCmd)
}

var RmIndexCmd = &cobra.Command{
	Use:   "rm_index PATHS...",
	Short: "Remove index entries",
	Long:  "Remove index entries",
	Run: func(c *cobra.Command, args []string) {
		err := Env.Init()
		if err != nil {
			Env.Fatalf("Unable to initialize environment: %s\n", err.Error())
		}
		defer Env.Close()

		controller.RmIndex(Env, args)
	},
}
