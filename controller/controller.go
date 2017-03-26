package controller

import (
	"errors"
	"fmt"
	"github.com/atongen/gosaic/environment"
	"github.com/atongen/gosaic/model"
	"github.com/atongen/gosaic/service"
	"github.com/atongen/gosaic/util"
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

func findOrCreateProject(env environment.Environment, inPath, name, coverOutfile, macroOutfile, mosaicOutfile string) (*model.Project, error) {
	if inPath == "" {
		return nil, errors.New("Error: path is empty")
	}

	if _, err := os.Stat(inPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("Error: file not found: %s\n", inPath)
	}

	dir, err := filepath.Abs(filepath.Dir(inPath))
	if err != nil {
		return nil, fmt.Errorf("Error getting image directory: %s\n", err.Error())
	}

	fName := filepath.Base(inPath)
	ext := filepath.Ext(fName)
	extL := strings.ToLower(ext)
	if !util.SliceContainsString([]string{".jpg", ".jpeg", ".png"}, extL) {
		return nil, errors.New("Error: only jpg images can be processed")
	}

	basename := fName[:len(fName)-len(ext)]
	if name == "" {
		name = basename
	}

	baseFilename := util.CleanStr(name)

	projectService := env.ServiceFactory().MustProjectService()
	project, err := projectService.GetOneBy("name = ?", name)
	if err != nil {
		return nil, fmt.Errorf("Error getting project name: %s\n", err.Error())
	}

	if project == nil {
		project = &model.Project{Name: name}
	} else {
		var status, action string
		if project.IsComplete {
			status = "Complete"
			action = "rebuild"
		} else {
			status = "Incomplete"
			action = "resume"
		}
		fmt.Printf("%s project with name '%s' already exists. Type 'Y' to %s, or anything else to abort [Y,n]: ", status, name, action)
		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "Y" {
			return nil, fmt.Errorf("Not attempting to %s project.", action)
		}
	}
	project.Path = inPath

	if coverOutfile == "" {
		coverOutfile, err = util.NextAvailableFilename(filepath.Join(dir, baseFilename+"-cover.png"))
		if err != nil {
			return nil, fmt.Errorf("Error getting next available filename for cover: %s\n", err.Error())
		}
	}
	project.CoverPath = coverOutfile

	if macroOutfile == "" {
		macroOutfile, err = util.NextAvailableFilename(filepath.Join(dir, baseFilename+"-macro"+ext))
		if err != nil {
			return nil, fmt.Errorf("Error getting next available filename for macro: %s\n", err.Error())
		}
	}
	project.MacroPath = macroOutfile

	if mosaicOutfile == "" {
		mosaicOutfile, err = util.NextAvailableFilename(filepath.Join(dir, baseFilename+"-mosaic"+ext))
		if err != nil {
			return nil, fmt.Errorf("Error getting next available filename for mosaic: %s\n", err.Error())
		}
	}
	project.MosaicPath = mosaicOutfile

	if project.Id == int64(0) {
		err = projectService.Insert(project)
	} else {
		_, err = projectService.Update(project)
	}

	if err != nil {
		return nil, fmt.Errorf("Error creating project: %s\n", err.Error())
	}

	return project, nil
}

func envProject(env environment.Environment) (*model.Project, error) {
	if env.ProjectId() == int64(0) {
		return nil, nil
	}

	projectService := env.ServiceFactory().MustProjectService()
	return projectService.Get(env.ProjectId())
}

func envCover(env environment.Environment) (*model.Cover, error) {
	project, err := envProject(env)
	if err != nil {
		return nil, err
	}

	if project == nil {
		return nil, nil
	}

	if project.CoverId == int64(0) {
		return nil, nil
	}

	coverService := env.ServiceFactory().MustCoverService()
	return coverService.Get(project.CoverId)
}

func setEnvCover(env environment.Environment, cover *model.Cover) error {
	project, err := envProject(env)
	if err != nil {
		return err
	}

	if project == nil {
		return nil
	}

	projectService := env.ServiceFactory().MustProjectService()
	project.CoverId = cover.Id
	_, err = projectService.Update(project)
	return err
}

func envMacro(env environment.Environment) (*model.Macro, error) {
	project, err := envProject(env)
	if err != nil {
		return nil, err
	}

	if project == nil {
		return nil, nil
	}

	if project.MacroId == int64(0) {
		return nil, nil
	}

	macroService := env.ServiceFactory().MustMacroService()
	return macroService.Get(project.MacroId)
}

func setEnvMacro(env environment.Environment, macro *model.Macro) error {
	project, err := envProject(env)
	if err != nil {
		return err
	}

	if project == nil {
		return nil
	}

	projectService := env.ServiceFactory().MustProjectService()
	project.MacroId = macro.Id
	_, err = projectService.Update(project)
	return err
}

func envMosaic(env environment.Environment) (*model.Mosaic, error) {
	project, err := envProject(env)
	if err != nil {
		return nil, err
	}

	if project == nil {
		return nil, nil
	}

	if project.MosaicId == int64(0) {
		return nil, nil
	}

	mosaicService := env.ServiceFactory().MustMosaicService()
	return mosaicService.Get(project.MosaicId)
}

func setEnvMosaic(env environment.Environment, mosaic *model.Mosaic) error {
	project, err := envProject(env)
	if err != nil {
		return err
	}

	if project == nil {
		return nil
	}

	projectService := env.ServiceFactory().MustProjectService()
	project.MosaicId = mosaic.Id
	_, err = projectService.Update(project)
	return err
}

func projectComplete(env environment.Environment, project *model.Project) error {
	if project.IsComplete {
		return nil
	}

	projectService := env.ServiceFactory().MustProjectService()
	project.IsComplete = true
	_, err := projectService.Update(project)
	return err
}

// projectCleanup deletes all partial comparisons for a macro
func projectCleanup(env environment.Environment, macro *model.Macro) error {
	partialComparisonService := env.ServiceFactory().MustPartialComparisonService()
	return partialComparisonService.DeleteFrom(macro)
}
