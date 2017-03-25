package cmd

import (
	"fmt"
	"os"

	"github.com/atongen/gosaic/controller"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(StatusCmd)
}

var StatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show environment status",
	Long:  "Show environment status",
	Run: func(c *cobra.Command, args []string) {
		err := Env.Init()
		if err != nil {
			fmt.Printf("Unable to initialize environment: %s\n", err.Error())
			os.Exit(1)
		}
		defer Env.Close()

		controller.Status(Env)
	},
}
