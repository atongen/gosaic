package controller

import (
	"errors"
	"github.com/atongen/gosaic/environment"
	"github.com/atongen/gosaic/model"
	"github.com/atongen/gosaic/util"
	"image"

	"github.com/disintegration/imaging"

	"gopkg.in/cheggaaa/pb.v1"
)

func Macro(env environment.Environment, path string, coverId int64, outfile string) *model.Macro {
	coverService := env.ServiceFactory().MustCoverService()

	cover, err := coverService.Get(coverId)
	if err != nil {
		env.Printf("Error getting cover: %s\n", err.Error())
		return nil
	} else if cover == nil {
		env.Printf("Cover id %d not found\n", coverId)
		return nil
	}

	macro, img, err := findOrCreateMacro(env, cover, path, outfile)
	if err != nil {
		env.Printf("Error creating macro: %s\n", err.Error())
		return nil
	}

	err = buildMacroPartials(env, img, macro, env.Workers())
	if err != nil {
		env.Printf("Error creating macro partials: %s\n", err.Error())
		return nil
	}

	return macro
}

func findOrCreateMacro(env environment.Environment, cover *model.Cover, path, outfile string) (*model.Macro, *image.Image, error) {
	macroService := env.ServiceFactory().MustMacroService()
	aspectService := env.ServiceFactory().MustAspectService()

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
		env.Printf("Wrote macro image: %s\n", outfile)
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

func buildMacroPartials(env environment.Environment, img *image.Image, macro *model.Macro, workers int) error {
	macroPartialService := env.ServiceFactory().MustMacroPartialService()

	countMissing, err := macroPartialService.CountMissing(macro)
	if err != nil {
		return err
	}

	if countMissing == 0 {
		return nil
	}

	env.Printf("Building %d macro partials...\n", countMissing)

	bar := pb.StartNew(int(countMissing))

	batchSize := 100
	errs := make(chan error)

	go func(myErrs <-chan error) {
		for e := range myErrs {
			env.Printf("Error building macro partial: %s\n", e.Error())
		}
	}(errs)

	for {
		if env.Cancel() {
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

		processMacroPartials(env, img, macro, coverPartials, errs, workers)
		bar.Add(num)
	}

	close(errs)
	bar.Finish()
	return err
}

func processMacroPartials(env environment.Environment, img *image.Image, macro *model.Macro, coverPartials []*model.CoverPartial, errs chan<- error, workers int) {
	add := make(chan *model.MacroPartial)
	sem := make(chan bool, workers)

	go storeMacroPartials(env, add, sem, errs)

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

func storeMacroPartials(env environment.Environment, add <-chan *model.MacroPartial, sem <-chan bool, errs chan<- error) {
	macroPartialService := env.ServiceFactory().MustMacroPartialService()

	for macroPartial := range add {
		err := macroPartialService.Insert(macroPartial)
		if err != nil {
			errs <- err
		}
		<-sem
	}
}
