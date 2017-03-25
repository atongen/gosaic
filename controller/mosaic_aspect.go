package controller

import (
	"github.com/atongen/gosaic/environment"
	"github.com/atongen/gosaic/model"
)

func MosaicAspect(env environment.Environment,
	inPath, name, fillType string,
	coverWidth, coverHeight, partialWidth, partialHeight, size, maxRepeats int,
	threashold float64,
	coverOutfile, macroOutfile, mosaicOutfile string,
	cleanup, destructive bool) *model.Mosaic {

	project, err := findOrCreateProject(env, inPath, name, coverOutfile, macroOutfile, mosaicOutfile)
	if err != nil {
		env.Println(err.Error())
		return nil
	}
	env.SetProjectId(project.Id)

	cover, macro := MacroAspect(env, project.Path, coverWidth, coverHeight, partialWidth, partialHeight, size, project.CoverPath, project.MacroPath)
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

	mosaic := MosaicBuild(env, fillType, macro.Id, maxRepeats, destructive)
	if mosaic == nil {
		return nil
	}

	err = MosaicDraw(env, mosaic.Id, project.MosaicPath)
	if err != nil {
		return nil
	}

	err = projectComplete(env, project)
	if err != nil {
		return nil
	}

	if cleanup {
		err = projectCleanup(env, macro)
		if err != nil {
			return nil
		}
	}

	return mosaic
}
