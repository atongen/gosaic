package controller

import (
	"gosaic/environment"
	"gosaic/model"
	"gosaic/service"
	"gosaic/util"
	"image"
	"log"

	"github.com/disintegration/imaging"

	"gopkg.in/cheggaaa/pb.v1"
)

func Macro(env environment.Environment, path string, coverId int64, outfile string) *model.Macro {
	aspectService, err := env.AspectService()
	if err != nil {
		env.Printf("Error getting aspect service: %s\n", err.Error())
		return nil
	}

	coverService, err := env.CoverService()
	if err != nil {
		env.Printf("Error creating cover service: %s\n", err.Error())
		return nil
	}

	macroService, err := env.MacroService()
	if err != nil {
		env.Printf("Error creating macro service: %s\n", err.Error())
		return nil
	}

	macroPartialService, err := env.MacroPartialService()
	if err != nil {
		env.Printf("Error creating macro partial service: %s\n", err.Error())
		return nil
	}

	cover, err := coverService.Get(coverId)
	if err != nil {
		env.Printf("Error getting cover: %s\n", err.Error())
		return nil
	} else if cover == nil {
		env.Printf("Cover id %d not found\n", coverId)
		return nil
	}

	aspect, err := aspectService.Get(cover.AspectId)
	if err != nil {
		env.Printf("Error getting cover aspect: %s\n", err.Error())
		return nil
	}

	md5sum, err := util.Md5sum(path)
	if err != nil {
		env.Printf("Error getting macro md5sum: %s\n", err.Error())
		return nil
	}

	img, err := util.OpenImage(path)
	if err != nil {
		env.Printf("Failed to open image: %s\n", err.Error())
		return nil
	}

	orientation, err := util.GetOrientation(path)
	if err != nil {
		env.Printf("Failed to get image orientation: %s\n", err.Error())
		return nil
	}

	err = util.FixOrientation(img, orientation)
	if err != nil {
		env.Printf("Failed to fix image orientation: %s\n", err.Error())
		return nil
	}
	bounds := (*img).Bounds()

	var imgCov image.Image
	imgCov = imaging.Fill(*img, int(cover.Width), int(cover.Height), imaging.Center, imaging.Lanczos)

	if outfile != "" {
		env.Printf("Wrote resized macro image to %s\n", outfile)
		err = imaging.Save(imgCov, outfile)
		if err != nil {
			env.Printf("Error saving file: %s\n", err.Error())
			return nil
		}
	}

	macro, _ := macroService.GetOneBy("cover_id = ? AND md5sum = ?", cover.Id, md5sum)

	if macro == nil {
		macro = &model.Macro{
			AspectId:    aspect.Id,
			CoverId:     cover.Id,
			Path:        path,
			Md5sum:      md5sum,
			Width:       uint(bounds.Max.X),
			Height:      uint(bounds.Max.Y),
			Orientation: orientation,
		}
		err = macroService.Insert(macro)
		if err != nil {
			env.Printf("Error creating macro: %s\n", err.Error())
			return nil
		}
	}

	err = buildMacroPartials(env.Log(), macroPartialService, &imgCov, macro, env.Workers())
	if err != nil {
		env.Printf("Error building macro partials: %s\n", err.Error())
		return nil
	}

	env.Printf("Created macro for path %s with cover %d\n", path, cover.Id)

	return macro
}

func buildMacroPartials(l *log.Logger, macroPartialService service.MacroPartialService, img *image.Image, macro *model.Macro, workers int) error {
	countMissing, err := macroPartialService.CountMissing(macro)
	if err != nil {
		return nil
	}

	if countMissing == 0 {
		return nil
	}

	l.Printf("Building %d macro partials...\n", countMissing)

	bar := pb.StartNew(int(countMissing))

	batchSize := 100
	errs := make(chan error)

	go func(myLog *log.Logger, myErrs <-chan error) {
		for e := range myErrs {
			l.Printf("Error building macro partial: %s\n", e.Error())
		}
	}(l, errs)

	for {
		var coverPartials []*model.CoverPartial

		coverPartials, err = macroPartialService.FindMissing(macro, "cover_partials.id ASC", batchSize, 0)
		if err != nil {
			return err
		}

		num := len(coverPartials)
		if num == 0 {
			break
		}

		processMacroPartials(macroPartialService, img, macro, coverPartials, errs, workers)
		bar.Add(num)
	}

	bar.Finish()

	close(errs)

	return err
}

func processMacroPartials(macroPartialService service.MacroPartialService, img *image.Image, macro *model.Macro, coverPartials []*model.CoverPartial, errs chan<- error, workers int) {
	add := make(chan *model.MacroPartial)
	sem := make(chan bool, workers)

	go storeMacroPartials(macroPartialService, add, sem, errs)

	for _, coverPartial := range coverPartials {
		sem <- true
		go storeMacroPartial(img, macro, coverPartial, add, errs)
	}

	for i := 0; i < cap(sem); i++ {
		sem <- true
	}

	close(add)
	close(sem)
}

func storeMacroPartial(img *image.Image, macro *model.Macro, coverPartial *model.CoverPartial, add chan<- *model.MacroPartial, errs chan<- error) {
	macroPartial := model.MacroPartial{
		MacroId:        macro.Id,
		CoverPartialId: coverPartial.Id,
		AspectId:       coverPartial.AspectId,
	}

	pixels, err := util.GetImgPartialLab(img, coverPartial)
	if err != nil {
		errs <- err
		return
	}
	macroPartial.Pixels = pixels

	add <- &macroPartial
}

func storeMacroPartials(macroPartialService service.MacroPartialService, add <-chan *model.MacroPartial, sem <-chan bool, errs chan<- error) {
	for macroPartial := range add {
		err := macroPartialService.Insert(macroPartial)
		if err != nil {
			errs <- err
		}
		<-sem
	}
}
