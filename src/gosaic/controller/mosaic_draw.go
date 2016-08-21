package controller

import (
	"gosaic/environment"
	"gosaic/model"
	"gosaic/service"
	"gosaic/util"
	"image/color"
	"log"

	"gopkg.in/cheggaaa/pb.v1"

	"github.com/disintegration/imaging"
)

func MosaicDraw(env environment.Environment, mosaicId int64, outfile string) {
	macroService, err := env.MacroService()
	if err != nil {
		env.Printf("Error getting macro service: %s\n", err.Error())
		return
	}

	coverService, err := env.CoverService()
	if err != nil {
		env.Printf("Error getting cover service: %s\n", err.Error())
		return
	}

	mosaicService, err := env.MosaicService()
	if err != nil {
		env.Printf("Error getting mosaic service: %s\n", err.Error())
		return
	}

	mosaicPartialService, err := env.MosaicPartialService()
	if err != nil {
		env.Printf("Error getting mosaic partial service: %s\n", err.Error())
		return
	}

	mosaic, err := mosaicService.Get(mosaicId)
	if err != nil {
		env.Printf("Error getting mosaic id %d: %s\n", mosaicId, err.Error())
		return
	}

	if mosaic == nil {
		env.Printf("Mosaic id %d does not exist", mosaicId)
		return
	}

	macro, err := macroService.Get(mosaic.MacroId)
	if err != nil {
		env.Printf("Error getting macro: %s\n", err.Error())
		return
	}

	if macro == nil {
		env.Printf("Macro id %d does not exist", mosaic.MacroId)
		return
	}

	cover, err := coverService.Get(macro.CoverId)
	if err != nil {
		env.Printf("Error getting cover: %s\n", err.Error())
		return
	}

	if cover == nil {
		env.Printf("cover id %d does not exist", macro.CoverId)
		return
	}

	err = drawMosaic(env.Log(), mosaic, cover, mosaicPartialService, outfile)
	if err != nil {
		env.Printf("Error drawing mosaic: %s\n", err.Error())
	}
	env.Printf("Wrote mosaic %s to %s\n", mosaic.Name, outfile)
}

func drawMosaic(l *log.Logger, mosaic *model.Mosaic, cover *model.Cover, mosaicPartialService service.MosaicPartialService, outfile string) error {
	numPartials, err := mosaicPartialService.Count(mosaic)
	if err != nil {
		return err
	}

	if numPartials == 0 {
		l.Println("This mosaic has 0 partials")
		return nil
	}

	dst := imaging.New(int(cover.Width), int(cover.Height), color.NRGBA{0, 0, 0, 0})

	batchSize := 100
	numCreated := 0

	l.Printf("Drawing %d mosaic partials\n", numPartials)
	bar := pb.StartNew(int(numPartials))

	for {
		mosaicPartialViews, err := mosaicPartialService.FindAllPartialViews(mosaic, "mosaic_partials.id asc", batchSize, numCreated)
		if err != nil {
			return err
		}

		num := len(mosaicPartialViews)
		if num == 0 {
			break
		}

		for _, view := range mosaicPartialViews {
			img, err := util.GetImageCoverPartial(view.Gidx, view.CoverPartial)
			if err != nil {
				return err
			}
			dst = imaging.Paste(dst, *img, view.CoverPartial.Pt())
			bar.Increment()
		}

		numCreated += num
	}

	bar.Finish()

	err = imaging.Save(dst, outfile)
	if err != nil {
		return err
	}

	return nil
}
