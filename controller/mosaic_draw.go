package controller

import (
	"errors"
	"fmt"
	"github.com/atongen/gosaic/environment"
	"github.com/atongen/gosaic/model"
	"github.com/atongen/gosaic/util"
	"image/color"

	"gopkg.in/cheggaaa/pb.v1"

	"github.com/disintegration/imaging"
)

func MosaicDraw(env environment.Environment, mosaicId int64, outfile string) error {
	macroService := env.ServiceFactory().MustMacroService()
	coverService := env.ServiceFactory().MustCoverService()
	mosaicService := env.ServiceFactory().MustMosaicService()

	mosaic, err := mosaicService.Get(mosaicId)
	if err != nil {
		env.Printf("Error getting mosaic id %d: %s\n", mosaicId, err.Error())
		return err
	}

	if mosaic == nil {
		msg := fmt.Sprintf("Mosaic id %d does not exist\n", mosaicId)
		env.Println(msg)
		return errors.New(msg)
	}

	macro, err := macroService.Get(mosaic.MacroId)
	if err != nil {
		env.Printf("Error getting macro: %s\n", err.Error())
		return err
	}

	if macro == nil {
		msg := fmt.Sprintf("Macro id %d does not exist\n", mosaic.MacroId)
		env.Println(msg)
		return errors.New(msg)
	}

	cover, err := coverService.Get(macro.CoverId)
	if err != nil {
		env.Printf("Error getting cover: %s\n", err.Error())
		return err
	}

	if cover == nil {
		msg := fmt.Sprintf("Cover id %d does not exist\n", macro.CoverId)
		env.Println(msg)
		return errors.New(msg)
	}

	err = drawMosaic(env, mosaic, cover, outfile)
	if err != nil {
		env.Printf("Error drawing mosaic: %s\n", err.Error())
		return err
	}
	env.Printf("Wrote mosaic image: %s\n", outfile)

	writeExif("", macro.Path, outfile)

	return nil
}

func drawMosaic(env environment.Environment, mosaic *model.Mosaic, cover *model.Cover, outfile string) error {
	mosaicPartialService := env.ServiceFactory().MustMosaicPartialService()

	numPartials, err := mosaicPartialService.Count(mosaic)
	if err != nil {
		return err
	}

	if numPartials == 0 {
		env.Println("This mosaic has 0 partials")
		return nil
	}

	dst := imaging.New(int(cover.Width), int(cover.Height), color.NRGBA{0, 0, 0, 0})

	batchSize := 100
	numCreated := 0

	env.Printf("Drawing %d mosaic partials...\n", numPartials)
	bar := pb.StartNew(int(numPartials))

	for {
		if env.Cancel() {
			return errors.New("Cancelled")
		}

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

func writeExif(toolPath, src, dst string) error {
	tp, err := util.ExiftoolPath(toolPath)
	if err != nil {
		return err
	} else if tp == "" {
		return nil
	}
	_, err = util.ExifCp(tp, src, dst)
	return err
}
