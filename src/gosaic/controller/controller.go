package controller

import (
	"errors"
	"fmt"
	"gosaic/model"
	"gosaic/service"
	"gosaic/util"
	"os"
	"path/filepath"
	"strings"
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

func validateMosaicArgs(mosaicService service.MosaicService, inPath, name, coverOutfile, macroOutfile, mosaicOutfile string) (string, string, string, string, error) {
	if inPath == "" {
		return "", "", "", "", errors.New("path is empty")
	}

	if _, err := os.Stat(inPath); os.IsNotExist(err) {
		return "", "", "", "", fmt.Errorf("file not found: %s\n", inPath)
	}

	dir, err := filepath.Abs(filepath.Dir(inPath))
	if err != nil {
		return "", "", "", "", fmt.Errorf("Error getting image directory: %s\n", err.Error())
	}

	fName := filepath.Base(inPath)
	ext := filepath.Ext(fName)
	extL := strings.ToLower(ext)
	if extL != ".jpg" && extL != ".jpeg" {
		return "", "", "", "", errors.New("only jpg images can be processed")
	}

	basename := fName[:len(fName)-len(ext)]
	if name == "" {
		name = basename
	}

	found, err := mosaicService.ExistsBy("name = ?", name)
	if err != nil {
		return "", "", "", "", fmt.Errorf("Error checking for mosaic name uniqueness: %s\n", err.Error())
	} else if found {
		return "", "", "", "", fmt.Errorf("Mosaic with name '%s' already exists\n", name)
	}

	baseFilename := util.CleanStr(name)

	if coverOutfile == "" {
		coverOutfile = filepath.Join(dir, baseFilename+"-cover.png")
	}
	if _, err := os.Stat(coverOutfile); err == nil {
		return "", "", "", "", fmt.Errorf("cover out file already exists: %s\n", coverOutfile)
	}

	if macroOutfile == "" {
		macroOutfile = filepath.Join(dir, baseFilename+"-macro"+ext)
	}
	if _, err := os.Stat(macroOutfile); err == nil {
		return "", "", "", "", fmt.Errorf("macro out file already exists: %s\n", macroOutfile)
	}

	if mosaicOutfile == "" {
		mosaicOutfile = filepath.Join(dir, baseFilename+"-mosaic"+ext)
	}
	if _, err := os.Stat(macroOutfile); err == nil {
		return "", "", "", "", fmt.Errorf("mosaic out file already exists: %s\n", mosaicOutfile)
	}

	return name, coverOutfile, macroOutfile, mosaicOutfile, nil
}
