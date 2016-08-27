package controller

import (
	"errors"
	"gosaic/model"
	"gosaic/service"
	"gosaic/util"
)

func getImageAspect(path string, aspectService service.AspectService) (*model.Aspect, error) {
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

func calculateDimensions(aspectService service.AspectService, path string, width, height int) (int, int, error) {
	if width < 0 || height < 0 {
		return 0, 0, errors.New("Width and height must be greater than or equal to zero")
	}

	if width > 0 && height > 0 {
		return width, height, nil
	}

	aspect, err := getImageAspect(path, aspectService)
	if err != nil {
		return 0, 0, err
	}

	var cWidth, cHeight int

	if width == 0 {
		cWidth = aspect.RoundWidth(height)
	} else {
		cWidth = width
	}

	if height == 0 {
		cHeight = aspect.RoundHeight(width)
	} else {
		cHeight = height
	}

	return cWidth, cHeight, nil
}
