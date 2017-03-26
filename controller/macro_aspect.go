package controller

import (
	"github.com/atongen/gosaic/environment"
	"github.com/atongen/gosaic/model"
)

func MacroAspect(env environment.Environment, path string, coverWidth, coverHeight, partialWidth, partialHeight, size int, coverOutfile, macroOutfile string) (*model.Cover, *model.Macro) {
	aspectService := env.ServiceFactory().MustAspectService()

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

	cover, err := envCover(env)
	if err != nil {
		env.Printf("Error getting cover from project environment: %s\n", err.Error())
		return nil, nil
	}

	if cover == nil {
		cover = CoverAspect(env, myCoverWidth, myCoverHeight, myPartialWidth, myPartialHeight, size)
		if cover == nil {
			env.Println("Failed to create cover")
			return nil, nil
		}
	}

	err = setEnvCover(env, cover)
	if err != nil {
		env.Printf("Error setting cover in project environment: %s\n", err.Error())
		return nil, nil
	}

	if coverOutfile != "" {
		err = CoverDraw(env, cover.Id, coverOutfile)
		if err != nil {
			env.Printf("Error drawing cover: %s\n", err.Error())
			return cover, nil
		}
	}

	macro, err := envMacro(env)
	if err != nil {
		env.Printf("Error getting macro from project environment: %s\n", err.Error())
		return cover, nil
	}

	if macro == nil {
		macro = Macro(env, path, cover.Id, macroOutfile)
		if macro == nil {
			env.Println("Failed to create macro")
			return cover, nil
		}
	}

	err = setEnvMacro(env, macro)
	if err != nil {
		env.Printf("Error setting macro in project environment: %s\n", err.Error())
		return cover, nil
	}

	return cover, macro
}
