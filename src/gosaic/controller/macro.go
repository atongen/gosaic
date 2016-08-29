package controller

import (
	"errors"
	"gosaic/environment"
	"gosaic/model"
	"gosaic/service"
	"gosaic/util"
	"image"
	"log"
	"os"
	"os/signal"
	"syscall"

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

	macro, img, err := findOrCreateMacro(macroService, aspectService, cover, path, outfile)
	if err != nil {
		env.Printf("Error creating macro: %s\n", err.Error())
		return nil
	}

	err = buildMacroPartials(env.Log(), macroPartialService, img, macro, env.Workers())
	if err != nil {
		env.Printf("Error creating macro partials: %s\n", err.Error())
		return nil
	}

	return macro
}

func findOrCreateMacro(macroService service.MacroService, aspectService service.AspectService, cover *model.Cover, path, outfile string) (*model.Macro, *image.Image, error) {
	md5sum, err := util.Md5sum(path)
	if err != nil {
		return nil, nil, err
	}

	img, err := util.OpenImage(path)
	if err != nil {
		return nil, nil, err
	}

	orientation, err := util.GetOrientation(path)
	if err != nil {
		return nil, nil, err
	}

	err = util.FixOrientation(img, orientation)
	if err != nil {
		return nil, nil, err
	}
	bounds := (*img).Bounds()

	var imgCov image.Image
	imgCov = imaging.Fill(*img, cover.Width, cover.Height, imaging.Center, imaging.Lanczos)

	if outfile != "" {
		err = imaging.Save(imgCov, outfile)
		if err != nil {
			return nil, nil, err
		}
	}

	macro, err := macroService.GetOneBy("cover_id = ? AND md5sum = ?", cover.Id, md5sum)
	if err != nil {
		return nil, nil, err
	}

	if macro == nil {
		aspect, err := aspectService.FindOrCreate(bounds.Max.X, bounds.Max.Y)
		if err != nil {
			return nil, nil, err
		}

		macro = &model.Macro{
			AspectId:    aspect.Id,
			CoverId:     cover.Id,
			Path:        path,
			Md5sum:      md5sum,
			Width:       bounds.Max.X,
			Height:      bounds.Max.Y,
			Orientation: orientation,
		}
		err = macroService.Insert(macro)
		if err != nil {
			return nil, nil, err
		}
	}

	return macro, &imgCov, nil
}

func buildMacroPartials(l *log.Logger, macroPartialService service.MacroPartialService, img *image.Image, macro *model.Macro, workers int) error {
	countMissing, err := macroPartialService.CountMissing(macro)
	if err != nil {
		return err
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

	cancel := false
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cancel = true
	}()

	for {
		if cancel {
			return errors.New("Cancelled")
		}

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
		go storeMacroPartial(img, macro, coverPartial, add)
	}

	for i := 0; i < cap(sem); i++ {
		sem <- true
	}

	close(add)
	close(sem)
}

func storeMacroPartial(img *image.Image, macro *model.Macro, coverPartial *model.CoverPartial, add chan<- *model.MacroPartial) {
	macroPartial := model.MacroPartial{
		MacroId:        macro.Id,
		CoverPartialId: coverPartial.Id,
		AspectId:       coverPartial.AspectId,
	}

	pixels := util.GetImgPartialLab(img, coverPartial)
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
