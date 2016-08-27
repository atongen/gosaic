package controller

import (
	"gosaic/environment"
	"gosaic/model"
)

func MacroAspect(env environment.Environment, path string, coverWidth, coverHeight, partialWidth, partialHeight, num int, outfile string) (*model.Cover, *model.Macro) {
	aspectService, err := env.AspectService()
	if err != nil {
		env.Printf("Error getting aspect service: %s\n", err.Error())
		return nil, nil
	}

	myCoverWidth, myCoverHeight, err := calculateDimensions(aspectService, path, coverWidth, coverHeight)
	if err != nil {
		env.Printf("Error getting cover dimensions: %s\n", err.Error())
		return nil, nil
	}

	cover := CoverAspect(env, myCoverWidth, myCoverHeight, partialWidth, partialHeight, num)
	if cover == nil {
		env.Println("Failed to create cover")
		return nil, nil
	}
	macro := Macro(env, path, cover.Id, outfile)
	if macro == nil {
		env.Println("Failed to create macro")
		return cover, nil
	}

	return cover, macro
}
