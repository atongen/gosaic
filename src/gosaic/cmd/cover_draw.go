package cmd

import (
	"gosaic/controller"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	coverDrawOutPath string
)

func init() {
	addLocalFlag(&coverDrawOutPath, "out", "o", "", "Path to write output file", CoverDrawCmd)
	RootCmd.AddCommand(CoverDrawCmd)
}

var CoverDrawCmd = &cobra.Command{
	Use:   "cover_draw NAME",
	Short: "Draw a cover image",
	Long:  "Draw a cover image",
	Run: func(c *cobra.Command, args []string) {
		if len(args) != 1 {
			Env.Fatalln("name is required")
		}

		if coverDrawOutPath == "" {
			Env.Fatalln("out path is required")
		}

		err := Env.Init()
		if err != nil {
			Env.Fatalf("Unable to initialize environment: %s\n", err.Error())
		}
		defer Env.Close()

		ext := strings.ToLower(filepath.Ext(coverDrawOutPath))
		if ext != ".png" {
			Env.Fatalf("Out path must be a .png file")
		}

		controller.CoverDraw(Env, args[0], coverDrawOutPath)
	},
}
