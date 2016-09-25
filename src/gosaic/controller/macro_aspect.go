package controller

import (
	"gosaic/environment"
	"gosaic/model"
)

func MacroAspect(env environment.Environment, path string, coverWidth, coverHeight, partialWidth, partialHeight, num int, coverOutfile, macroOutfile string) (*model.Cover, *model.Macro) {
	aspectService := env.MustAspectService()

	aspect, width, height, err := getImageDimensions(aspectService, path)
	if err != nil {
		env.Printf("Error getting image aspect: %s\n", err.Error())
		return nil, nil
	}

	myCoverWidth, myCoverHeight := calculateDimensionsFromAspect(aspect, coverWidth, coverHeight, width, height)

	coverAspect, err := aspectService.FindOrCreate(myCoverWidth, myCoverHeight)
	if err != nil {
		env.Printf("Error getting cover aspect: %s\n", err.Error())
		return nil, nil
	}

	myPartialWidth, myPartialHeight := calculateDimensionsFromAspect(coverAspect, partialWidth, partialHeight, myCoverWidth, myCoverHeight)

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
