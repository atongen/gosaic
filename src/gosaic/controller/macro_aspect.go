package controller

import (
	"fmt"
	"gosaic/environment"
	"gosaic/model"
	"gosaic/service"
	"gosaic/util"
	"time"
)

func MacroAspect(env environment.Environment, path string, coverWidth, coverHeight, partialWidth, partialHeight, num int, outfile string) (*model.Cover, *model.Macro) {
	ts := time.Now().Format(time.RubyDate)
	name := fmt.Sprintf("%s-%s", path, ts)

	var myCoverWidth, myCoverHeight int

	if coverWidth < 0 || coverHeight < 0 {
		env.Println("Cover width and height must not be less than zero")
		return nil, nil
	}

	if coverWidth > 0 && coverHeight > 0 {
		myCoverWidth = coverWidth
		myCoverHeight = coverHeight
	} else {
		aspectService, err := env.AspectService()
		if err != nil {
			env.Printf("Error getting aspect service: %s\n", err.Error())
			return nil, nil
		}

		aspect, err := macroAspectGetImageAspect(path, aspectService)
		if err != nil {
			env.Printf("Error getting aspect: %s\n", err.Error())
			return nil, nil
		}

		if coverWidth == 0 {
			myCoverWidth = aspect.RoundWidth(coverHeight)
		} else {
			myCoverWidth = coverWidth
		}

		if coverHeight == 0 {
			myCoverHeight = aspect.RoundHeight(coverWidth)
		} else {
			myCoverHeight = coverHeight
		}
	}

	cover := CoverAspect(env, name, myCoverWidth, myCoverHeight, partialWidth, partialHeight, num)
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
