package controller

import (
	"gosaic/model"
	"gosaic/service"
	"gosaic/util"
)

func getImageDimensions(aspectService service.AspectService, path string) (*model.Aspect, int, int, error) {
	img, err := util.OpenImage(path)
	if err != nil {
		return nil, 0, 0, err
	}

	orientation, err := util.GetOrientation(path)
	if err != nil {
		return nil, 0, 0, err
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
		return nil, 0, 0, err
	}

	return aspect, width, height, nil
}

func calculateDimensions(aspectService service.AspectService, path string, targetWidth, targetHeight int) (int, int, error) {
	if targetWidth > 0 && targetHeight > 0 {
		return targetWidth, targetHeight, nil
	}

	aspect, width, height, err := getImageDimensions(aspectService, path)
	if err != nil {
		return 0, 0, err
	}

	cWidth, cHeight := calculateDimensionsFromAspect(aspect, targetWidth, targetHeight, width, height)
	return cWidth, cHeight, nil
}

func calculateDimensionsFromAspect(aspect *model.Aspect, width, height, baseWidth, baseHeight int) (int, int) {
	if width > 0 && height > 0 {
		return width, height
	} else if width <= 0 && height <= 0 {
		return baseWidth, baseHeight
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

	return cWidth, cHeight
}
