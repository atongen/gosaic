package cmd

import (
	"gosaic/controller"

	"github.com/spf13/cobra"
)

var (
	indexPath string
)

func init() {
	addLocalFlag(&indexPath, "index", "i", "", "Path to index", IndexCmd)
	RootCmd.AddCommand(IndexCmd)
}

var IndexCmd = &cobra.Command{
	Use:   "index",
	Short: "Add path to index",
	Long:  "Add path to index",
	Run: func(c *cobra.Command, args []string) {
		if indexPath == "" {
			Env.Fatalln("index path is required")
		}

		err := Env.Init()
		if err != nil {
			Env.Fatalf("Unable to initialize environment: %s\n", err.Error())
		}
		defer Env.Close()

		controller.Index(Env, indexPath)
	},
}
