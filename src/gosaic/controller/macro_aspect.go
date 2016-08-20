package controller

import (
	"fmt"
	"gosaic/environment"
	"gosaic/model"
	"gosaic/service"
	"gosaic/util"
	"time"
)

func MacroAspect(env environment.Environment, path string, coverWidth, coverHeight, partialWidth, partialHeight, num int) (*model.Cover, *model.Macro) {
	ts := time.Now().Format(time.RubyDate)
	name := fmt.Sprintf("%s-%s", path, ts)

	var myCoverWidth, myCoverHeight int

	if coverWidth == 0 && coverHeight == 0 {
		env.Fatalf("either width or height can be zero, but not both")
	} else if coverWidth == 0 || coverHeight == 0 {
		aspectService, err := env.AspectService()
		if err != nil {
			env.Fatalf("Error getting aspect service: %s\n", err.Error())
		}

		aspect, err := macroAspectGetImageAspect(path, aspectService)
		if err != nil {
			env.Fatalf("Error getting aspect: %s\n", err.Error())
		}

		if coverWidth == 0 {
			myCoverWidth = aspect.RoundWidth(coverHeight)
			myCoverHeight = coverHeight
		} else if coverHeight == 0 {
			myCoverWidth = coverWidth
			myCoverHeight = aspect.RoundHeight(coverWidth)
		}
	} else {
		myCoverWidth = coverWidth
		myCoverHeight = coverHeight
	}

	cover := CoverAspect(env, name, myCoverWidth, myCoverHeight, partialWidth, partialHeight, num)
	macro := Macro(env, path, cover.Id)

	return cover, macro
}

func macroAspectGetImageAspect(path string, aspectService service.AspectService) (*model.Aspect, error) {
	img, err := util.OpenImage(path)
	if err != nil {
		return nil, err
	}

	orientation, err := util.GetOrientation(path)
	if err != nil {
		return nil, err
	}

	swap := false
	if 4 < orientation && orientation <= 8 {
		swap = true
	}
	if orientation == 0 {
		orientation = 1
	}

	bounds := (*img).Bounds()

	var width, height int
	if swap {
		width = bounds.Max.Y
		height = bounds.Max.X
	} else {
		width = bounds.Max.X
		height = bounds.Max.Y
	}

	aspect, err := aspectService.FindOrCreate(width, height)
	if err != nil {
		return nil, err
	}

	return aspect, nil
}
