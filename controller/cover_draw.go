package controller

import (
	"errors"
	"github.com/atongen/gosaic/environment"
	"github.com/atongen/gosaic/model"

	"github.com/fogleman/gg"
)

func CoverDraw(env environment.Environment, coverId int64, outPath string) error {
	coverService := env.ServiceFactory().MustCoverService()

	cover, err := coverService.Get(coverId)
	if err != nil {
		return err
	} else if cover == nil {
		return errors.New("Cover not found")
	}

	err = doCoverDraw(env, cover, outPath)
	if err != nil {
		return err
	}
	env.Printf("Wrote cover image: %s\n", outPath)

	return nil
}

func doCoverDraw(env environment.Environment, cover *model.Cover, outPath string) error {
	coverPartialService := env.ServiceFactory().MustCoverPartialService()

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
