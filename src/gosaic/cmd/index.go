package cmd

import (
	"gosaic/controller"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(IndexCmd)
}

var IndexCmd = &cobra.Command{
	Use:   "index PATH",
	Short: "Add path to index",
	Long:  "Add path to index",
	Run: func(c *cobra.Command, args []string) {
		if len(args) == 0 {
			Env.Fatalln("index path is required")
		}

		err := Env.Init()
		if err != nil {
			Env.Fatalf("Unable to initialize environment: %s\n", err.Error())
		}
		defer Env.Close()

		controller.Index(Env, args[0])
	},
}
