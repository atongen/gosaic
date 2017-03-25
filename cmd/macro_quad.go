package cmd

import (
	"github.com/atongen/gosaic/controller"

	"github.com/spf13/cobra"
)

var (
	macroQuadWidth        int
	macroQuadHeight       int
	macroQuadSize         int
	macroQuadMinDepth     int
	macroQuadMaxDepth     int
	macroQuadMinArea      int
	macroQuadMaxArea      int
	macroQuadCoverOutfile string
	macroQuadMacroOutfile string
)

func init() {
	addLocalIntFlag(&macroQuadWidth, "width", "w", 0, "Pixel width of cover, 0 maintains aspect from height", MacroQuadCmd)
	addLocalIntFlag(&macroQuadHeight, "height", "", 0, "Pixel height of cover, 0 maintains aspect from width", MacroQuadCmd)
	addLocalIntFlag(&macroQuadSize, "size", "s", -1, "Number of times to subdivide the image into quads", MacroQuadCmd)
	addLocalIntFlag(&macroQuadMinDepth, "min-depth", "", -1, "Minimum depth of quad subdivisions", MacroQuadCmd)
	addLocalIntFlag(&macroQuadMaxDepth, "max-depth", "", -1, "Maximum depth of quad subdivisions", MacroQuadCmd)
	addLocalIntFlag(&macroQuadMinArea, "min-area", "", -1, "Minimum area of quad subdivisions", MacroQuadCmd)
	addLocalIntFlag(&macroQuadMinArea, "max-area", "", -1, "Maxumum area of quad subdivisions", MacroQuadCmd)
	addLocalStrFlag(&macroQuadCoverOutfile, "cover-out", "", "", "File to write cover image", MacroQuadCmd)
	addLocalStrFlag(&macroQuadMacroOutfile, "out", "o", "", "File to write resized macro image", MacroQuadCmd)
	RootCmd.AddCommand(MacroQuadCmd)
}

var MacroQuadCmd = &cobra.Command{
	Use:    "macro_quad PATH",
	Short:  "Add quad cover and macro",
	Long:   "Add quad cover and macro",
	Hidden: true,
	Run: func(c *cobra.Command, args []string) {
		if len(args) != 1 {
			Env.Fatalln("PATH is required")
		}

		if args[0] == "" {
			Env.Fatalln("Macro path is required")
		}

		if macroQuadWidth < 0 {
			Env.Fatalln("width must be greater than zero")
		}

		if macroQuadHeight < 0 {
			Env.Fatalln("height must be greater than zero")
		}

		if macroQuadSize == 0 &&
			macroQuadMinDepth == 0 &&
			macroQuadMaxDepth == 0 &&
			macroQuadMinArea == 0 &&
			macroQuadMaxArea == 0 {
			Env.Fatalln("Add least one of size, min-depth, max-depth, min-area, or max-area must be non-zero.")
		}

		err := Env.Init()
		if err != nil {
			Env.Fatalf("Unable to initialize environment: %s\n", err.Error())
		}
		defer Env.Close()

		controller.MacroQuad(
			Env,
			args[0],
			macroQuadWidth,
			macroQuadHeight,
			macroQuadSize,
			macroQuadMinDepth,
			macroQuadMaxDepth,
			macroQuadMinArea,
			macroQuadMaxArea,
			macroQuadCoverOutfile,
			macroQuadMacroOutfile,
		)
	},
}
