package controller

import (
	"gosaic/environment"
	"gosaic/model"
	"os"
	"path/filepath"
	"strings"
)

func MosaicQuad(env environment.Environment,
	inPath, name, fillType string,
	coverWidth, coverHeight, num, maxDepth, minArea, maxRepeats int,
	coverOutfile, macroOutfile, mosaicOutfile string) *model.Mosaic {

	if inPath == "" {
		env.Println("Error: path is empty")
		return nil
	}

	if _, err := os.Stat(inPath); os.IsNotExist(err) {
		env.Printf("Error: file not found: %s\n", inPath)
		return nil
	}

	dir, err := filepath.Abs(filepath.Dir(inPath))
	if err != nil {
		env.Printf("Error getting image directory: %s\n", err.Error())
		return nil
	}

	fName := filepath.Base(inPath)
	ext := filepath.Ext(fName)
	extL := strings.ToLower(ext)
	if extL != ".jpg" && extL != ".jpeg" {
		env.Println("Error: only jpg images can be processed")
		return nil
	}

	if fillType != "best" && fillType != "random" {
		env.Printf("Invalid mosaic fill type: %s\n", fillType)
		return nil
	}

	basename := fName[:len(fName)-len(ext)]
	if name == "" {
		name = basename
	}

	if coverOutfile == "" {
		coverOutfile = filepath.Join(dir, basename+"-cover"+ext)
	}

	if macroOutfile == "" {
		macroOutfile = filepath.Join(dir, basename+"-macro"+ext)
	}

	if mosaicOutfile == "" {
		mosaicOutfile = filepath.Join(dir, basename+"-mosaic"+ext)
	}

	cover, macro := MacroQuad(env, inPath, coverWidth, coverHeight, num, maxDepth, minArea, coverOutfile, macroOutfile)
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

	mosaic := MosaicBuild(env, name, fillType, macro.Id, maxRepeats)
	if mosaic == nil {
		return nil
	}

	err = MosaicDraw(env, mosaic.Id, mosaicOutfile)
	if err != nil {
		return nil
	}

	return mosaic
}
