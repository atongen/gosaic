package cmd

import (
	"fmt"

	"github.com/atongen/gosaic/environment"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(VersionCmd)
}

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number and exit",
	Long:  "Print the version number and exit",
	Run: func(c *cobra.Command, args []string) {
		fmt.Printf("gosaic %s (%s) %s %s\n",
			environment.Version,
			environment.BuildTime,
			environment.BuildUser,
			environment.BuildHash)
	},
}
