package controller

import (
	"gosaic/environment"
	"gosaic/util"
	"image/color"

	"github.com/disintegration/imaging"
)

func MosaicDraw(env environment.Environment, mosaicId int64, outfile string) {
	macroService, err := env.MacroService()
	if err != nil {
		env.Fatalf("Error getting macro service: %s\n", err.Error())
	}

	coverService, err := env.CoverService()
	if err != nil {
		env.Fatalf("Error getting cover service: %s\n", err.Error())
	}

	mosaicService, err := env.MosaicService()
	if err != nil {
		env.Fatalf("Error getting mosaic service: %s\n", err.Error())
	}

	mosaicPartialService, err := env.MosaicPartialService()
	if err != nil {
		env.Fatalf("Error getting mosaic partial service: %s\n", err.Error())
	}

	mosaic, err := mosaicService.Get(mosaicId)
	if err != nil {
		env.Fatalf("Error getting mosaic id %d: %s\n", mosaicId, err.Error())
	}

	if mosaic == nil {
		env.Fatalf("Mosaic id %d does not exist", mosaicId)
	}

	macro, err := macroService.Get(mosaic.MacroId)
	if err != nil {
		env.Fatalf("Error getting macro: %s\n", err.Error())
	}

	if macro == nil {
		env.Fatalf("Macro id %d does not exist", mosaic.MacroId)
	}

	cover, err := coverService.Get(macro.CoverId)
	if err != nil {
		env.Fatalf("Error getting cover: %s\n", err.Error())
	}

	if cover == nil {
		env.Fatalf("cover id %d does not exist", macro.CoverId)
	}

	numPartials, err := mosaicPartialService.Count(mosaic)
	if err != nil {
		env.Fatalf("Error counting mosaic partials: %s\n", err.Error())
	}

	if numPartials == 0 {
		env.Println("This mosaic has 0 partials")
		return
	}

	dst := imaging.New(int(cover.Width), int(cover.Height), color.NRGBA{0, 0, 0, 0})

	batchSize := 100
	numCreated := 0

	for {
		mosaicPartialViews, err := mosaicPartialService.FindAllPartialViews(mosaic, "mosaic_partials.id asc", batchSize, numCreated)
		if err != nil {
			env.Fatalf("Error finding mosaic partials: %s\n", err.Error())
		}

		if len(mosaicPartialViews) == 0 {
			break
		}

		for _, view := range mosaicPartialViews {
			img, err := util.GetImageCoverPartial(view.Gidx, view.CoverPartial)
			if err != nil {
				env.Fatalf("Error getting mosaic partial image: %s\n", err.Error())
			}
			dst = imaging.Paste(dst, *img, view.CoverPartial.Pt())
		}

		numCreated += len(mosaicPartialViews)
	}

	err = imaging.Save(dst, outfile)
	if err != nil {
		env.Fatalf("Error writing mosaic to %s: %s\n", outfile, err.Error())
	}

	env.Println("Wrote mosaic %s to %s\n", mosaic.Name, outfile)
}
