package cmd

import "github.com/spf13/cobra"

func init() {
	RootCmd.AddCommand(MosaicCmd)
}

var MosaicCmd = &cobra.Command{
	Use:   "mosaic",
	Short: "Create a mosaic",
	Long:  "Create a mosaic",
}
