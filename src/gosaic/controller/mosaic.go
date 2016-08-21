package controller

import (
	"gosaic/environment"
	"gosaic/model"
)

func Mosaic(env environment.Environment,
	path, name, mosaicType string,
	coverWidth, coverHeight, partialWidth, partialHeight, num, maxRepeats int,
	mosaicOutfile, macroOutfile string) *model.Mosaic {

	if mosaicType != "best" && mosaicType != "random" {
		env.Printf("Invalid mosaic build type: %s\n", mosaicType)
		return nil
	}

	cover, macro := MacroAspect(env, path, coverWidth, coverHeight, partialWidth, partialHeight, num, macroOutfile)
	if cover == nil || macro == nil {
		env.Printf("Failed to create cover or macro")
		return nil
	}

	PartialAspect(env, macro.Id)
	Compare(env, macro.Id)

	mosaic := MosaicBuild(env, name, mosaicType, macro.Id, maxRepeats)
	if mosaic == nil {
		env.Println("Failed to build mosaic")
		return nil
	}

	MosaicDraw(env, mosaic.Id, mosaicOutfile)

	return mosaic
}
