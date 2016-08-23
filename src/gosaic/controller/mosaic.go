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
		return nil
	}

	err := PartialAspect(env, macro.Id)
	if err != nil {
		return nil
	}

	err = Compare(env, macro.Id)
	if err != nil {
		return nil
	}

	mosaic := MosaicBuild(env, name, mosaicType, macro.Id, maxRepeats)
	if mosaic == nil {
		return nil
	}

	err = MosaicDraw(env, mosaic.Id, mosaicOutfile)
	if err != nil {
		return nil
	}

	return mosaic
}
