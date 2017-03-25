package cmd

import (
	"strconv"
	"strings"

	"github.com/atongen/gosaic/controller"
	"github.com/spf13/cobra"
)

var (
	mosaicAspectName          string
	mosaicAspectFillType      string
	mosaicAspectCoverWidth    int
	mosaicAspectCoverHeight   int
	mosaicAspectPartialAspect string
	mosaicAspectSize          int
	mosaicAspectMaxRepeats    int
	mosaicAspectThreashold    float64
	mosaicAspectOutfile       string
	mosaicAspectCoverOutfile  string
	mosaicAspectMacroOutfile  string
	mosaicAspectCleanup       bool
	mosaicAspectDestructive   bool
)

func init() {
	addLocalStrFlag(&mosaicAspectName, "name", "n", "", "Name of mosaic", MosaicAspectCmd)
	addLocalStrFlag(&mosaicAspectFillType, "fill-type", "f", "random", "Mosaic fill to use, either 'random' or 'best'", MosaicAspectCmd)
	addLocalIntFlag(&mosaicAspectCoverWidth, "width", "w", 0, "Pixel width of mosaic, 0 maintains aspect from image height", MosaicAspectCmd)
	addLocalIntFlag(&mosaicAspectCoverHeight, "height", "", 0, "Pixel height of mosaic, 0 maintains aspect from width", MosaicAspectCmd)
	addLocalStrFlag(&mosaicAspectPartialAspect, "aspect", "a", "", "Aspect of mosaic partials (CxR)", MosaicAspectCmd)
	addLocalIntFlag(&mosaicAspectSize, "size", "s", 0, "Number of mosaic partials in smallest dimension, 0 auto-calculates", MosaicAspectCmd)
	addLocalIntFlag(&mosaicAspectMaxRepeats, "max-repeats", "", -1, "Number of times an index image can be repeated, 0 is unlimited, -1 is the minimun number", MosaicAspectCmd)
	addLocalFloatFlag(&mosaicAspectThreashold, "threashold", "t", -1.0, "How similar aspect ratios must be", MosaicAspectCmd)
	addLocalStrFlag(&mosaicAspectOutfile, "out", "", "", "File to write final mosaic image", MosaicAspectCmd)
	addLocalStrFlag(&mosaicAspectCoverOutfile, "cover-out", "", "", "File to write cover partial pattern image", MosaicAspectCmd)
	addLocalStrFlag(&mosaicAspectMacroOutfile, "macro-out", "", "", "File to write resized macro image", MosaicAspectCmd)
	addLocalBoolFlag(&mosaicAspectCleanup, "cleanup", "", false, "Delete mosaic metadata after completion", MosaicAspectCmd)
	addLocalBoolFlag(&mosaicAspectDestructive, "destructive", "d", false, "Delete mosaic metadata during creation", MosaicAspectCmd)
	MosaicCmd.AddCommand(MosaicAspectCmd)
}

var MosaicAspectCmd = &cobra.Command{
	Use:   "aspect PATH",
	Short: "Create an aspect mosaic from image at PATH",
	Long:  "Create an aspect mosaic from image at PATH",
	Run: func(c *cobra.Command, args []string) {
		if len(args) != 1 {
			Env.Fatalln("Mosaic path is required")
		}

		if args[0] == "" {
			Env.Fatalln("Mosaic path is required")
		}

		if mosaicAspectCoverWidth < 0 {
			Env.Fatalln("width must be greater than zero")
		}

		if mosaicAspectCoverHeight < 0 {
			Env.Fatalln("height must be greater than zero")
		}

		var (
			aw, ah int
			err    error
		)

		if mosaicAspectPartialAspect == "" {
			aw = 0
			ah = 0
		} else {
			aspectStrings := strings.Split(mosaicAspectPartialAspect, "x")
			if len(aspectStrings) != 2 {
				Env.Fatalln("aspect format must be CxR")
			}

			aw, err = strconv.Atoi(aspectStrings[0])
			if err != nil {
				Env.Fatalf("Error converting aspect columns: %s\n", err.Error())
			}

			if aw < 0 {
				Env.Fatalln("aspect columns must be greater than zero")
			}

			ah, err = strconv.Atoi(aspectStrings[1])
			if err != nil {
				Env.Fatalf("Error converting aspect rows: %s\n", err.Error())
			}

			if ah < 0 {
				Env.Fatalln("aspect rows must be greater than zero")
			}
		}

		if mosaicAspectFillType != "best" && mosaicAspectFillType != "random" {
			Env.Fatalln("Invalid fill-type")
		}

		err = Env.Init()
		if err != nil {
			Env.Fatalf("Unable to initialize environment: %s\n", err.Error())
		}
		defer Env.Close()

		controller.MosaicAspect(
			Env,
			args[0],
			mosaicAspectName,
			mosaicAspectFillType,
			mosaicAspectCoverWidth,
			mosaicAspectCoverHeight,
			aw,
			ah,
			mosaicAspectSize,
			mosaicAspectMaxRepeats,
			mosaicAspectThreashold,
			mosaicAspectCoverOutfile,
			mosaicAspectMacroOutfile,
			mosaicAspectOutfile,
			mosaicAspectCleanup,
			mosaicAspectDestructive,
		)
	},
}
