package cmd

import (
	"gosaic/controller"

	"github.com/spf13/cobra"
)

var (
	coverSquareWidth  int
	coverSquareHeight int
	coverSquareNum    int
)

func init() {
	addLocalIntFlag(&coverSquareWidth, "width", "", 0, "Pixel width of cover", CoverSquareCmd)
	addLocalIntFlag(&coverSquareHeight, "height", "", 0, "Pixel height of cover", CoverSquareCmd)
	addLocalIntFlag(&coverSquareNum, "size", "s", 0, "Number of partials in smallest dimension", CoverSquareCmd)
	RootCmd.AddCommand(CoverSquareCmd)
}

var CoverSquareCmd = &cobra.Command{
	Use:   "cover_square NAME",
	Short: "Create a square cover",
	Long:  "Create a square cover",
	Run: func(c *cobra.Command, args []string) {
		if len(args) != 1 {
			Env.Fatalln("name is required")
		}

		if coverSquareWidth == 0 {
			Env.Fatalln("width is required")
		}

		if coverSquareHeight == 0 {
			Env.Fatalln("height is required")
		}

		if coverSquareNum == 0 {
			Env.Fatalln("num is required")
		}

		err := Env.Init()
		if err != nil {
			Env.Fatalf("Unable to initialize environment: %s\n", err.Error())
		}
		defer Env.Close()

		controller.CoverSquare(Env, args[0], coverSquareWidth, coverSquareHeight, coverSquareNum)
	},
}
