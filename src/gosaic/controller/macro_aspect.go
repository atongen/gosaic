package controller

import (
	"gosaic/environment"
	"gosaic/model"
)

func MacroAspect(env environment.Environment, path string, coverWidth, coverHeight, partialWidth, partialHeight, num int, coverOutfile, macroOutfile string) (*model.Cover, *model.Macro) {
	aspectService, err := env.AspectService()
	if err != nil {
		env.Printf("Error getting aspect service: %s\n", err.Error())
		return nil, nil
	}

	aspect, err := getImageAspect(path, aspectService)
	if err != nil {
		env.Printf("Error getting image aspect: %s\n", err.Error())
		return nil, nil
	}

	myCoverWidth, myCoverHeight, err := calculateDimensionsFromAspect(aspect, coverWidth, coverHeight)
	if err != nil {
		env.Printf("Error getting cover dimensions: %s\n", err.Error())
		return nil, nil
	}

	myPartialWidth, myPartialHeight, err := calculateDimensionsFromAspect(aspect, partialWidth, partialHeight)
	if err != nil {
		env.Printf("Error getting partial dimensions: %s\n", err.Error())
		return nil, nil
	}

	cover := CoverAspect(env, myCoverWidth, myCoverHeight, myPartialWidth, myPartialHeight, num)
	if cover == nil {
		env.Println("Failed to create cover")
		return nil, nil
	}

	if coverOutfile != "" {
		err = CoverDraw(env, cover.Id, coverOutfile)
		if err != nil {
			env.Printf("Error drawing cover: %s\n", err.Error())
			return cover, nil
		}
	}

	macro := Macro(env, path, cover.Id, macroOutfile)
	if macro == nil {
		env.Println("Failed to create macro")
		return cover, nil
	}

	return cover, macro
}
