package cmd

import (
	"gosaic/controller"

	"github.com/spf13/cobra"
)

var (
	macroQuadWidth   int
	macroQuadHeight  int
	macroQuadNum     int
	macroQuadOutfile string
)

func init() {
	addLocalIntFlag(&macroQuadWidth, "width", "", 0, "Pixel width of cover, 0 maintains aspect from height", MacroQuadCmd)
	addLocalIntFlag(&macroQuadHeight, "height", "", 0, "Pixel height of cover, 0 maintains aspect from width", MacroQuadCmd)
	addLocalIntFlag(&macroQuadNum, "num", "n", 0, "Number of times to subdivide the image into quads", MacroQuadCmd)
	addLocalFlag(&macroQuadOutfile, "out", "", "", "File to write resized macro image", MacroQuadCmd)
	RootCmd.AddCommand(MacroQuadCmd)
}

var MacroQuadCmd = &cobra.Command{
	Use:   "macro_quad PATH",
	Short: "Add quad cover and macro",
	Long:  "Add quad cover and macro",
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

		if macroQuadNum == 0 {
			Env.Fatalln("num is required")
		} else if macroQuadNum < 0 {
			Env.Fatalln("num must be greater than zero")
		}

		err := Env.Init()
		if err != nil {
			Env.Fatalf("Unable to initialize environment: %s\n", err.Error())
		}
		defer Env.Close()

		controller.MacroQuad(Env, args[0], macroQuadWidth, macroQuadHeight, macroQuadNum, macroQuadOutfile)
	},
}
