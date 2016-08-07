package controller

import (
	"gosaic/environment"
	"gosaic/model"
	"gosaic/service"
	"gosaic/util"
	"image"
)

func Macro(env environment.Environment, path, coverName string) {
	aspectService, err := env.AspectService()
	if err != nil {
		env.Printf("Error getting aspect service: %s\n", err.Error())
		return
	}

	coverService, err := env.CoverService()
	if err != nil {
		env.Printf("Error creating cover service: %s\n", err.Error())
		return
	}

	macroService, err := env.MacroService()
	if err != nil {
		env.Printf("Error creating macro service: %s\n", err.Error())
		return
	}

	img, err := util.OpenImage(path)
	if err != nil {
		env.Printf("Failed to open image: %s\n", err.Error())
		return
	}

	orientation, err := util.GetOrientation(path)
	if err != nil {
		env.Printf("Failed to get image orientation: %s\n", err.Error())
		return
	}

	err = util.FixOrientation(img, orientation)
	if err != nil {
		env.Printf("Failed to fix image orientation: %s\n", err.Error())
		return
	}

	bounds := (*img).Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y

	aspect, err := aspectService.FindOrCreate(width, height)
	if err != nil {
		env.Println("Error getting image aspect data", path, err)
		return
	}

	cover, err := coverService.GetOneBy("name", coverName)
	if err != nil {
		env.Printf("Error getting cover: %s\n", err.Error())
		return
	} else if cover == nil {
		env.Printf("Cover %s not found\n", coverName)
		return
	}

	coverAspect, err := aspectService.Get(cover.AspectId)
	if err != nil {
		env.Printf("Error getting cover aspect: %s\n", err.Error())
		return
	} else if cover == nil {
		env.Println("Cover aspect not found")
		return
	}

	if aspect.Id != coverAspect.Id {
		env.Printf("Aspect of image (%dx%d) does not match aspect of cover %s (%dx%d)\n",
			aspect.Columns, aspect.Rows, cover.Name, coverAspect.Columns, coverAspect.Rows)
		return
	}

	md5sum, err := util.Md5sum(path)
	if err != nil {
		env.Printf("Error getting macro md5sum: %s\n", err.Error())
		return
	}

	var macro *model.Macro
	macro, _ = macroService.GetOneBy("cover_id = ? AND md5sum = ?", cover.Id, md5sum)

	// macro was not found
	if macro.Id == int64(0) {
		macro = &model.Macro{
			AspectId:    aspect.Id,
			CoverId:     cover.Id,
			Path:        path,
			Md5sum:      md5sum,
			Width:       uint(width),
			Height:      uint(height),
			Orientation: orientation,
		}
		err = macroService.Insert(macro)
		if err != nil {
			env.Printf("Error creating macro: %s\n", err.Error())
			return
		}
	}

	err = buildMacroPartials(env, img, macro, env.Workers())
	if err != nil {
		env.Printf("Error building macro partials: %s\n", err.Error())
	}
}

func buildMacroPartials(env environment.Environment, img *image.Image, macro *model.Macro, workers int) error {
	macroPartialService, err := env.MacroPartialService()
	if err != nil {
		return err
	}

	batchSize := 1000
	errs := make(chan error)

	go func(myErrs <-chan error) {
		for e := range myErrs {
			env.Println(e.Error())
		}
	}(errs)

	for i := 0; ; i++ {
		var coverPartials []*model.CoverPartial

		coverPartials, err = macroPartialService.FindMissing(macro, "cover_partials.id ASC", batchSize, 0)
		if err != nil {
			return err
		}

		num := len(coverPartials)
		if num == 0 {
			break
		}

		env.Printf("Processing %d macro partials\n", num)
		processMacroPartials(macroPartialService, img, macro, coverPartials, errs, workers)
		if err != nil {
			break
		}

	}

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

	pixels, err := util.GetImgPartialLab(img, macro, coverPartial)
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
