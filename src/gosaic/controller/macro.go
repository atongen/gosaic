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
		env.Fatalf("Error getting aspect service: %s\n", err.Error())
	}

	coverService, err := env.CoverService()
	if err != nil {
		env.Fatalf("Error creating cover service: %s\n", err.Error())
	}

	macroService, err := env.MacroService()
	if err != nil {
		env.Fatalf("Error creating macro service: %s\n", err.Error())
	}

	cover, err := coverService.GetOneBy("name", coverName)
	if err != nil {
		env.Fatalf("Error getting cover: %s\n", err.Error())
	} else if cover == nil {
		env.Fatalf("Cover %s not found\n", coverName)
	}

	aspect, err := aspectService.Get(cover.AspectId)
	if err != nil {
		env.Fatalf("Error getting cover aspect: %s\n", err.Error())
	}

	md5sum, err := util.Md5sum(path)
	if err != nil {
		env.Fatalf("Error getting macro md5sum: %s\n", err.Error())
	}

	img, err := util.OpenImage(path)
	if err != nil {
		env.Fatalf("Failed to open image: %s\n", err.Error())
	}

	orientation, err := util.GetOrientation(path)
	if err != nil {
		env.Fatalf("Failed to get image orientation: %s\n", err.Error())
	}

	err = util.FixOrientation(img, orientation)
	if err != nil {
		env.Fatalf("Failed to fix image orientation: %s\n", err.Error())
	}

	img = util.FillAspect(img, aspect)
	bounds := (*img).Bounds()
	// width and height of image after resize to fill cover aspect
	width := bounds.Max.X
	height := bounds.Max.Y

	checkAspect, err := aspectService.Find(width, height)
	if err != nil {
		env.Fatalf("Error checking aspect: %s\n", err.Error())
	}

	if checkAspect == nil {
		env.Fatalf("No aspect for resized image found")
	}

	if aspect.Id != checkAspect.Id {
		env.Fatalf("Aspect of image (%dx%d) does not match aspect of cover (%dx%d)\n",
			checkAspect.Columns, checkAspect.Rows, aspect.Columns, aspect.Rows)
	}

	macro, _ := macroService.GetOneBy("cover_id = ? AND md5sum = ?", cover.Id, md5sum)

	if macro == nil {
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
			env.Fatalf("Error creating macro: %s\n", err.Error())
		}
	}

	err = buildMacroPartials(env, img, macro, env.Workers())
	if err != nil {
		env.Fatalf("Error building macro partials: %s\n", err.Error())
	}

	env.Printf("Built macro for %s with cover %s\n", path, coverName)
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
