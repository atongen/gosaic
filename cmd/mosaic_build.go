package cmd

import (
	"github.com/atongen/gosaic/controller"
	"github.com/spf13/cobra"
)

var (
	mosaicBuildMaxRepeats  int
	mosaicBuildMacroId     int
	mosaicBuildFillType    string
	mosaicBuildDestructive bool
)

func init() {
	addLocalIntFlag(&mosaicBuildMacroId, "macro-id", "", 0, "Id of macro to use to build mosaic", MosaicBuildCmd)
	addLocalIntFlag(&mosaicBuildMaxRepeats, "max-repeats", "", -1, "Number of times an index image can be repeated in the mosaic, 0 indicates unlimited, -1 is the minimum number", MosaicBuildCmd)
	addLocalStrFlag(&mosaicBuildFillType, "fill-type", "f", "random", "Mosaic build type, either 'best' or 'random'", MosaicBuildCmd)
	addLocalBoolFlag(&mosaicBuildDestructive, "destructive", "d", false, "Delete mosaic metadata during creation", MosaicBuildCmd)
	RootCmd.AddCommand(MosaicBuildCmd)
}

var MosaicBuildCmd = &cobra.Command{
	Use:    "mosaic_build NAME",
	Short:  "Build mosaic",
	Long:   "Build mosaic",
	Hidden: true,
	Run: func(c *cobra.Command, args []string) {
		if mosaicBuildMacroId == 0 {
			Env.Fatalln("Macro id is required")
		}

		if mosaicBuildFillType != "best" && mosaicBuildFillType != "random" {
			Env.Fatalln("type must be either 'best' or 'random'")
		}

		err := Env.Init()
		if err != nil {
			Env.Fatalf("Unable to initialize environment: %s\n", err.Error())
		}
		defer Env.Close()

		controller.MosaicBuild(Env, mosaicBuildFillType, int64(mosaicBuildMacroId), mosaicBuildMaxRepeats, mosaicBuildDestructive)
	},
}
