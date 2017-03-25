package cmd

import (
	"github.com/atongen/gosaic/controller"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(MosaicListCmd)
}

var MosaicListCmd = &cobra.Command{
	Use:    "mosaic_list",
	Short:  "List mosaic entries",
	Long:   "List mosaic entries",
	Hidden: true,
	Run: func(c *cobra.Command, args []string) {
		err := Env.Init()
		if err != nil {
			Env.Fatalf("Unable to initialize environment: %s\n", err.Error())
		}
		defer Env.Close()

		controller.MosaicList(Env)
	},
}
