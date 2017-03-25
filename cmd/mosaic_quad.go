package cmd

import (
	"github.com/atongen/gosaic/controller"
	"github.com/spf13/cobra"
)

var (
	mosaicQuadName         string
	mosaicQuadFillType     string
	mosaicQuadCoverWidth   int
	mosaicQuadCoverHeight  int
	mosaicQuadSize         int
	mosaicQuadMinDepth     int
	mosaicQuadMaxDepth     int
	mosaicQuadMinArea      int
	mosaicQuadMaxArea      int
	mosaicQuadMaxRepeats   int
	mosaicQuadThreashold   float64
	mosaicQuadOutfile      string
	mosaicQuadCoverOutfile string
	mosaicQuadMacroOutfile string
	mosaicQuadCleanup      bool
	mosaicQuadDestructive  bool
)

func init() {
	addLocalStrFlag(&mosaicQuadName, "name", "n", "", "Name of mosaic", MosaicQuadCmd)
	addLocalStrFlag(&mosaicQuadFillType, "fill-type", "f", "random", "Mosaic fill to use, either 'random' or 'best'", MosaicQuadCmd)
	addLocalIntFlag(&mosaicQuadCoverWidth, "width", "w", 0, "Pixel width of mosaic, 0 maintains aspect from image height", MosaicQuadCmd)
	addLocalIntFlag(&mosaicQuadCoverHeight, "height", "", 0, "Pixel height of mosaic, 0 maintains aspect from width", MosaicQuadCmd)
	addLocalIntFlag(&mosaicQuadSize, "size", "s", -1, "Number of times to split the partials into quads", MosaicQuadCmd)
	addLocalIntFlag(&mosaicQuadMinDepth, "min-depth", "", -1, "Minimum number of times all partials will be split into quads", MosaicQuadCmd)
	addLocalIntFlag(&mosaicQuadMaxDepth, "max-depth", "", -1, "Number of times a partial can be split into quads", MosaicQuadCmd)
	addLocalIntFlag(&mosaicQuadMinArea, "min-area", "", -1, "The smallest a partial can get before it can't be split", MosaicQuadCmd)
	addLocalIntFlag(&mosaicQuadMaxArea, "max-area", "", -1, "The largest a partial can be", MosaicQuadCmd)
	addLocalIntFlag(&mosaicQuadMaxRepeats, "max-repeats", "", -1, "Number of times an index image can be repeated, 0 is unlimited, -1 is the minimun number", MosaicQuadCmd)
	addLocalFloatFlag(&mosaicQuadThreashold, "threashold", "t", -1.0, "How similar aspect ratios must be", MosaicQuadCmd)
	addLocalStrFlag(&mosaicQuadOutfile, "out", "o", "", "File to write final mosaic image", MosaicQuadCmd)
	addLocalStrFlag(&mosaicQuadCoverOutfile, "cover-out", "", "", "File to write cover partial pattern image", MosaicQuadCmd)
	addLocalStrFlag(&mosaicQuadMacroOutfile, "macro-out", "", "", "File to write resized macro image", MosaicQuadCmd)
	addLocalBoolFlag(&mosaicQuadCleanup, "cleanup", "", false, "Delete mosaic metadata after completion", MosaicQuadCmd)
	addLocalBoolFlag(&mosaicQuadDestructive, "destructive", "d", false, "Delete mosaic metadata during creation", MosaicQuadCmd)
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

		if mosaicQuadFillType != "best" && mosaicAspectFillType != "random" {
			Env.Fatalln("Invalid fill-type")
		}

		if mosaicQuadSize == 0 &&
			mosaicQuadMinDepth == 0 &&
			mosaicQuadMaxDepth == 0 &&
			mosaicQuadMinArea == 0 &&
			mosaicQuadMaxArea == 0 {
			Env.Fatalln("Add least one of size, min-depth, max-depth, min-area, or max-area must be non-zero.")
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
			mosaicQuadSize,
			mosaicQuadMinDepth,
			mosaicQuadMaxDepth,
			mosaicQuadMinArea,
			mosaicQuadMaxArea,
			mosaicQuadMaxRepeats,
			mosaicQuadThreashold,
			mosaicQuadCoverOutfile,
			mosaicQuadMacroOutfile,
			mosaicQuadOutfile,
			mosaicQuadCleanup,
			mosaicQuadDestructive,
		)
	},
}
