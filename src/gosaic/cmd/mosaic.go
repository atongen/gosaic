package cmd

import (
	"gosaic/controller"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var (
	mosaicName          string
	mosaicType          string
	mosaicCoverWidth    int
	mosaicCoverHeight   int
	mosaicPartialAspect string
	mosaicSize          int
	mosaicMaxRepeats    int
	mosaicOutfile       string
	mosaicMacroOutfile  string
)

func init() {
	addLocalFlag(&mosaicName, "name", "", "", "Name of mosaic", MosaicCmd)
	addLocalFlag(&mosaicType, "type", "", "random", "Type of mosaic to build, either 'best' or 'random'", MosaicCmd)
	addLocalIntFlag(&mosaicCoverWidth, "width", "", 0, "Pixel width of mosaic, 0 maintains aspect from image height", MosaicCmd)
	addLocalIntFlag(&mosaicCoverHeight, "height", "", 0, "Pixel height of mosaic, 0 maintains aspect from width", MosaicCmd)
	addLocalFlag(&mosaicPartialAspect, "aspect", "a", "1x1", "Aspect of mosaic partials (CxR)", MosaicCmd)
	addLocalIntFlag(&mosaicSize, "size", "s", 10, "Number of mosaic partials in smallest dimension", MosaicCmd)
	addLocalIntFlag(&mosaicMaxRepeats, "max_repeats", "", -1, "Number of times an index image can be repeated, 0 is unlimited, -1 is the minimun number", MosaicCmd)
	addLocalFlag(&mosaicOutfile, "out", "", "", "File to write final mosaic image", MosaicCmd)
	addLocalFlag(&mosaicMacroOutfile, "macro_out", "", "", "File to write resized macro image", MosaicCmd)
	RootCmd.AddCommand(MosaicCmd)
}

var MosaicCmd = &cobra.Command{
	Use:   "mosaic PATH",
	Short: "Create mosaic from image at PATH",
	Long:  "Create mosaic from image at PATH",
	Run: func(c *cobra.Command, args []string) {
		if len(args) != 1 {
			Env.Fatalln("Mosaic path is required")
		}

		if args[0] == "" {
			Env.Fatalln("Mosaic path is required")
		}

		if mosaicName == "" {
			Env.Fatalln("Mosaic name is required")
		}

		if mosaicCoverWidth < 0 {
			Env.Fatalln("width must be greater than zero")
		}

		if mosaicCoverHeight < 0 {
			Env.Fatalln("height must be greater than zero")
		}

		if mosaicPartialAspect == "" {
			Env.Fatalln("aspect is required")
		}

		aspectStrings := strings.Split(mosaicPartialAspect, "x")
		if len(aspectStrings) != 2 {
			Env.Fatalln("aspect format must be CxR")
		}

		aw, err := strconv.Atoi(aspectStrings[0])
		if err != nil {
			Env.Fatalf("Error converting aspect columns: %s\n", err.Error())
		}

		if aw == 0 {
			Env.Fatalln("aspect columns cannot be zero")
		} else if aw < 0 {
			Env.Fatalln("aspect columns must be greater than zero")
		}

		ah, err := strconv.Atoi(aspectStrings[1])
		if err != nil {
			Env.Fatalf("Error converting aspect rows: %s\n", err.Error())
		}

		if ah == 0 {
			Env.Fatalln("aspect rows cannot be zero")
		} else if ah < 0 {
			Env.Fatalln("aspect rows must be greater than zero")
		}

		if mosaicSize == 0 {
			Env.Fatalln("size is required")
		} else if mosaicSize < 0 {
			Env.Fatalln("size must be greater than zero")
		}

		if mosaicOutfile == "" {
			Env.Fatalln("Mosaic outfile is required")
		}

		err = Env.Init()
		if err != nil {
			Env.Fatalf("Unable to initialize environment: %s\n", err.Error())
		}
		defer Env.Close()

		controller.Mosaic(Env, args[0], mosaicName, mosaicType, mosaicCoverWidth, mosaicCoverHeight, aw, ah, mosaicSize, mosaicMaxRepeats, mosaicOutfile, mosaicMacroOutfile)
	},
}
