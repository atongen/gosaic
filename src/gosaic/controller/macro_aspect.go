package controller

import (
	"fmt"
	"gosaic/environment"
	"gosaic/model"
	"gosaic/util"
	"time"
)

func MacroAspect(env environment.Environment, path string, coverWidth, coverHeight, partialWidth, partialHeight, num int) {
	coverService, err := env.CoverService()
	if err != nil {
		env.Printf("Error creating cover service: %s\n", err.Error())
		return
	}

	coverPartialService, err := env.CoverPartialService()
	if err != nil {
		env.Printf("Error creating cover partial service: %s\n", err.Error())
		return
	}

	aspectService, err := env.AspectService()
	if err != nil {
		env.Printf("Error getting aspect service: %s\n", err.Error())
		return
	}

	macroService, err := env.MacroService()
	if err != nil {
		env.Printf("Error creating macro service: %s\n", err.Error())
		return
	}

	macroPartialService, err := env.MacroPartialService()
	if err != nil {
		env.Printf("Error creating macro partial service: %s\n", err.Error())
		return
	}

	md5sum, err := util.Md5sum(path)
	if err != nil {
		env.Printf("Error getting macro md5sum: %s\n", err.Error())
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

	img = util.FillAspect(img, coverWidth, coverHeight)
	bounds := (*img).Bounds()
	// width and height of image after resize to fill cover aspect
	width := bounds.Max.X
	height := bounds.Max.Y

	coverAspect, err := aspectService.FindOrCreate(width, height)
	if err != nil {
		env.Printf("Error getting cover aspect: %s\n", err.Error())
		return
	}

	coverPartialAspect, err := aspectService.FindOrCreate(partialWidth, partialHeight)
	if err != nil {
		env.Printf("Error getting cover partial aspect: %s\n", err.Error())
		return
	}

	ts := time.Now().Format(time.RubyDate)
	name := fmt.Sprintf("%s-%s", path, ts)

	cover, err := createCoverAspect(coverService, coverAspect, name, coverWidth, coverHeight)
	if err != nil {
		env.Printf("Error creating aspect cover: %s\n", err.Error())
		return
	}
	env.Printf("Created cover %s\n", cover.Name)

	numPartials, err := addCoverAspectPartials(coverPartialService, cover, coverPartialAspect, num)
	if err != nil {
		env.Printf("Error adding aspect cover partials: %s\n", err.Error())
		return
	}
	env.Printf("Added %d aspect cover partials\n", numPartials)

	macro := &model.Macro{
		AspectId:    coverAspect.Id,
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
	env.Printf("Created macro for %s\n", path)

	err = buildMacroPartials(env, img, macro, env.Workers())
	if err != nil {
		env.Printf("Error building macro partials: %s\n", err.Error())
		return
	}

	numMacroPartials, err := macroPartialService.CountBy("macro_id", macro.Id)
	if err != nil {
		env.Printf("Error counting macro partials: %s\n", err.Error())
		return
	}

	env.Printf("Built %d macro partials", numMacroPartials)
}
