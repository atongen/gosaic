package controller

import (
	"gosaic/environment"
	"gosaic/model"
	"gosaic/service"
	"gosaic/util"
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

	macroPartialService, err := env.MacroPartialService()
	if err != nil {
		env.Printf("Error creating macro partial service: %s\n", err.Error())
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
	width := bounds.Max.Y
	height := bounds.Max.X

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
	}

	if aspect.Id != cover.AspectId {
		env.Printf("Aspect of image (%dx%d) does not match aspect of cover %s\n",
			path, aspect.Columns, aspect.Rows, cover.Name)
		return
	}

	md5sum, err := util.Md5sum(path)
	if err != nil {
		env.Printf("Error getting macro md5sum: %s\n", err.Error())
		return
	}

	macro, err := macroService.GetOneBy("cover_id = ? AND md5sum = ?", cover.Id, md5sum)
	if err != nil {
		env.Printf("Error checking macro existence: %s\n", err.Error())
		return
	}

	if macro == nil {
		macro = model.Macro{
			AspectId:    aspect.Id,
			CoverId:     cover.Id,
			Path:        path,
			Md5sum:      md5sum,
			Width:       uint(width),
			Height:      uint(height),
			Orientation: orientation,
		}
		err = macroService.Insert(&macro)
		if err != nil {
			env.Printf("Error creating macro: %s\n", err.Error())
			return
		}
	}

	num, err := buildMacroPartials(macroPartialService, macro)
	if err != nil {
		env.Printf("Error building macro partials: %s\n", err.Error())
		return
	}

	env.Printf("Created macro and built %d partials\n", num)
}

func buildMacroPartials(macroPartialService service.MacroPartialService, macro *model.Macro) (int, error) {
	batchSize := 1000
	num := 0

	for {
		var coverPartials []*model.CoverPartial

		coverPartials, err = macroPartialService.FindMissing(macro, "cover_partials.id ASC", batchSize, 0)
		if err != nil {
			return nil, err
		}

		// create batch here

		if len(coverPartials) == 0 {
			break
		}
	}

	return num, nil

}
