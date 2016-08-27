package controller

import (
	"errors"
	"gosaic/environment"
	"gosaic/model"
	"gosaic/service"

	"github.com/fogleman/gg"
)

func CoverDraw(env environment.Environment, coverId int64, outPath string) error {
	coverService, err := env.CoverService()
	if err != nil {
		return err
	}

	coverPartialService, err := env.CoverPartialService()
	if err != nil {
		return err
	}

	cover, err := coverService.Get(coverId)
	if err != nil {
		return err
	} else if cover == nil {
		return errors.New("Cover not found")
	}

	err = doCoverDraw(cover, outPath, coverPartialService)
	if err != nil {
		return err
	}

	return nil
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
