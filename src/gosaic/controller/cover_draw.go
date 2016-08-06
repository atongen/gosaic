package controller

import (
	"gosaic/environment"
	"gosaic/model"
	"gosaic/service"

	"github.com/fogleman/gg"
)

func CoverDraw(env environment.Environment, name, outPath string) {
	coverService, err := env.CoverService()
	if err != nil {
		env.Printf("Error creating cover service: %s\n", err.Error())
		return
	}

	coverPartialService, err := env.CoverPartialService()
	if err != nil {
		env.Printf("Error creating cover partial service: %s\n", err.Error())
		return
	}

	cover, err := coverService.GetOneBy("name", name)
	if err != nil {
		env.Printf("Error getting cover: %s\n", err.Error())
		return
	} else if cover == nil {
		env.Printf("Cover %s not found\n", name)
	}

	err = doCoverDraw(cover, outPath, coverPartialService)
	if err != nil {
		env.Printf("Error drawing cover: %s\n", err.Error())
	}
}

func doCoverDraw(cover *model.Cover, outPath string, coverPartialService service.CoverPartialService) error {
	dc := gg.NewContext(int(cover.Width), int(cover.Height))
	dc.Clear()

	coverPartials, err := coverPartialService.FindAll(cover.Id, "id ASC")
	if err != nil {
		return err
	}

	for _, cp := range coverPartials {
		x := float64(cp.X1)
		y := float64(cp.Y1)
		w := float64(cp.X2 - cp.X1)
		h := float64(cp.Y2 - cp.Y1)

		dc.Push()
		dc.SetRGBA(1, 0, 0, 0.5)
		dc.DrawRectangle(x, y, w, h)
		dc.SetLineWidth(2)
		dc.Stroke()
		dc.Pop()
	}

	return dc.SavePNG(outPath)
}
