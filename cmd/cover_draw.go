package cmd

import (
	"path/filepath"
	"strings"

	"github.com/atongen/gosaic/controller"
	"github.com/spf13/cobra"
)

var (
	coverDrawId int
)

func init() {
	addLocalIntFlag(&coverDrawId, "cover-id", "", 0, "Id of cover to draw", CoverDrawCmd)
	RootCmd.AddCommand(CoverDrawCmd)
}

var CoverDrawCmd = &cobra.Command{
	Use:    "cover_draw PATH",
	Short:  "Draw a cover image",
	Long:   "Draw a cover image",
	Hidden: true,
	Run: func(c *cobra.Command, args []string) {
		if len(args) != 1 {
			Env.Fatalln("Path is required")
		}

		if coverDrawId == 0 {
			Env.Fatalln("cover id is required")
		}

		err := Env.Init()
		if err != nil {
			Env.Fatalf("Unable to initialize environment: %s\n", err.Error())
		}
		defer Env.Close()

		ext := strings.ToLower(filepath.Ext(args[0]))
		if ext != ".png" {
			Env.Fatalf("Out path must be a .png file")
		}

		controller.CoverDraw(Env, int64(coverDrawId), args[0])
	},
}
