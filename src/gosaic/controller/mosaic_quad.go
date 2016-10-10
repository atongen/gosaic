package controller

import (
	"gosaic/environment"
	"gosaic/model"
)

func MosaicQuad(env environment.Environment,
	inPath, name, fillType string,
	coverWidth, coverHeight, num, maxDepth, minArea, maxRepeats int,
	threashold float64,
	coverOutfile, macroOutfile, mosaicOutfile string) *model.Mosaic {

	project, err := findOrCreateProject(env, inPath, name, coverOutfile, macroOutfile, mosaicOutfile)
	if err != nil {
		env.Println(err.Error())
		return nil
	}
	env.SetProjectId(project.Id)

	cover, macro := MacroQuad(env, project.Path, coverWidth, coverHeight, num, maxDepth, minArea, project.CoverPath, project.MacroPath)
	if cover == nil || macro == nil {
		return nil
	}

	err = PartialAspect(env, macro.Id, threashold)
	if err != nil {
		return nil
	}

	err = Compare(env, macro.Id)
	if err != nil {
		return nil
	}

	mosaic := MosaicBuild(env, fillType, macro.Id, maxRepeats)
	if mosaic == nil {
		return nil
	}

	err = MosaicDraw(env, mosaic.Id, project.MosaicPath)
	if err != nil {
		return nil
	}

	return mosaic
}
