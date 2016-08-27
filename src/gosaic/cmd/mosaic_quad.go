package cmd

import (
	"gosaic/controller"

	"github.com/spf13/cobra"
)

var (
	mosaicQuadName         string
	mosaicQuadFillType     string
	mosaicQuadCoverWidth   int
	mosaicQuadCoverHeight  int
	mosaicQuadNum          int
	mosaicQuadMaxDepth     int
	mosaicQuadMinArea      int
	mosaicQuadMaxRepeats   int
	mosaicQuadOutfile      string
	mosaicQuadCoverOutfile string
	mosaicQuadMacroOutfile string
)

func init() {
	addLocalStrFlag(&mosaicQuadName, "name", "n", "", "Name of mosaic", MosaicQuadCmd)
	addLocalStrFlag(&mosaicQuadFillType, "fill-type", "f", "random", "Mosaic fill to use, either 'random' or 'best'", MosaicQuadCmd)
	addLocalIntFlag(&mosaicQuadCoverWidth, "width", "w", 0, "Pixel width of mosaic, 0 maintains aspect from image height", MosaicQuadCmd)
	addLocalIntFlag(&mosaicQuadCoverHeight, "height", "", 0, "Pixel height of mosaic, 0 maintains aspect from width", MosaicQuadCmd)
	addLocalIntFlag(&mosaicQuadNum, "num", "", 1024, "Number of times to split the partials into quads", MosaicQuadCmd)
	addLocalIntFlag(&mosaicQuadMaxDepth, "max-depth", "", 0, "Number of times a partial can be split into quads", MosaicQuadCmd)
	addLocalIntFlag(&mosaicQuadMinArea, "min-area", "", 0, "The smallest an partial can get before it can't be split", MosaicQuadCmd)
	addLocalStrFlag(&mosaicQuadOutfile, "out", "o", "", "File to write final mosaic image", MosaicQuadCmd)
	addLocalStrFlag(&mosaicQuadCoverOutfile, "cover-out", "", "", "File to write cover partial pattern image", MosaicQuadCmd)
	addLocalStrFlag(&mosaicQuadMacroOutfile, "macro-out", "", "", "File to write resized macro image", MosaicQuadCmd)
	MosaicCmd.AddCommand(MosaicQuadCmd)
}

var MosaicQuadCmd = &cobra.Command{
	Use:   "quad PATH",
	Short: "Create quad-tree mosaic from image at PATH",
	Long:  "Create quad-tree mosaic from image at PATH",
	Run: func(c *cobra.Command, args []string) {
		if len(args) != 1 {
			Env.Fatalln("Mosaic path is required")
		}

		if args[0] == "" {
			Env.Fatalln("Mosaic path is required")
		}

		if mosaicQuadCoverWidth < 0 {
			Env.Fatalln("width must be greater than zero")
		}

		if mosaicQuadCoverHeight < 0 {
			Env.Fatalln("height must be greater than zero")
		}

		if mosaicQuadNum <= 0 {
			Env.Fatalln("num is required and must be greater than zero")
		}

		err := Env.Init()
		if err != nil {
			Env.Fatalf("Unable to initialize environment: %s\n", err.Error())
		}
		defer Env.Close()

		controller.MosaicQuad(
			Env,
			args[0],
			mosaicQuadName,
			mosaicQuadFillType,
			mosaicQuadCoverWidth,
			mosaicQuadCoverHeight,
			mosaicQuadNum,
			mosaicQuadMaxDepth,
			mosaicQuadMinArea,
			mosaicQuadMaxRepeats,
			mosaicQuadCoverOutfile,
			mosaicQuadMacroOutfile,
			mosaicQuadOutfile,
		)
	},
}
