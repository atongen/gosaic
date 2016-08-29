package controller

import (
	"gosaic/environment"
	"gosaic/model"
)

func MosaicQuad(env environment.Environment,
	inPath, name, fillType string,
	coverWidth, coverHeight, num, maxDepth, minArea, maxRepeats int,
	coverOutfile, macroOutfile, mosaicOutfile string) *model.Mosaic {

	mosaicService, err := env.MosaicService()
	if err != nil {
		env.Printf("Error getting mosaic service: %s\n", err.Error())
		return nil
	}

	myName, myCoverOutfile, myMacroOutfile, myMosaicOutfile, err := validateMosaicArgs(
		mosaicService, inPath, name, coverOutfile, macroOutfile, mosaicOutfile,
	)
	if err != nil {
		env.Printf("Error validating mosaic arguments: %s\n", err.Error())
		return nil
	}

	cover, macro := MacroQuad(env, inPath, coverWidth, coverHeight, num, maxDepth, minArea, myCoverOutfile, myMacroOutfile)
	if cover == nil || macro == nil {
		return nil
	}

	err = PartialAspect(env, macro.Id)
	if err != nil {
		return nil
	}

	err = Compare(env, macro.Id)
	if err != nil {
		return nil
	}

	mosaic := MosaicBuild(env, myName, fillType, macro.Id, maxRepeats)
	if mosaic == nil {
		return nil
	}

	err = MosaicDraw(env, mosaic.Id, myMosaicOutfile)
	if err != nil {
		return nil
	}

	return mosaic
}
