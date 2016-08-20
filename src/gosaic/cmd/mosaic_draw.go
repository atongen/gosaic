package cmd

import (
	"gosaic/controller"

	"github.com/spf13/cobra"
)

var (
	mosaicId int
)

func init() {
	addLocalIntFlag(&mosaicId, "mosaic_id", "", 0, "Id of mosaic to draw", MosaicDrawCmd)
	RootCmd.AddCommand(MosaicDrawCmd)
}

var MosaicDrawCmd = &cobra.Command{
	Use:   "mosaic_draw OUTFILE",
	Short: "Draw mosaic",
	Long:  "Draw mosaic",
	Run: func(c *cobra.Command, args []string) {
		if len(args) != 1 {
			Env.Fatalln("Mosaic out file is required")
		}

		if args[0] == "" {
			Env.Fatalln("Mosaic out file is required")
		}

		if mosaicId == 0 {
			Env.Fatalln("Mosaic id is required")
		}

		err := Env.Init()
		if err != nil {
			Env.Fatalf("Unable to initialize environment: %s\n", err.Error())
		}
		defer Env.Close()

		controller.MosaicDraw(Env, int64(mosaicId), args[0])
	},
}
