package cmd

import (
	"github.com/atongen/gosaic/controller"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(CoverListCmd)
}

var CoverListCmd = &cobra.Command{
	Use:    "cover_list",
	Short:  "List cover entries",
	Long:   "List cover entries",
	Hidden: true,
	Run: func(c *cobra.Command, args []string) {
		err := Env.Init()
		if err != nil {
			Env.Fatalf("Unable to initialize environment: %s\n", err.Error())
		}
		defer Env.Close()

		controller.CoverList(Env)
	},
}
