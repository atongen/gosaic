package controller

import (
	"errors"
	"fmt"
	"gosaic/environment"
	"gosaic/model"
	"gosaic/util"
	"image"
	"math"

	"gopkg.in/cheggaaa/pb.v1"
)

func MacroQuad(env environment.Environment,
	path string,
	coverWidth, coverHeight, size, maxDepth, minArea int,
	coverOutfile, macroOutfile string) (*model.Cover, *model.Macro) {

	aspectService := env.MustAspectService()
	coverService := env.MustCoverService()

	myCoverWidth, myCoverHeight, err := calculateDimensions(aspectService, path, coverWidth, coverHeight)
	if err != nil {
		env.Printf("Error getting cover dimensions: %s\n", err.Error())
		return nil, nil
	}

	size, maxDepth, minArea = macroQuadFixArgs(myCoverWidth, myCoverHeight, size, maxDepth, minArea)

	cover, err := envCover(env)
	if err != nil {
		env.Printf("Error getting cover from project environment: %s\n", err.Error())
		return nil, nil
	}

	if cover == nil {
		cover, err = macroQuadCreateCover(env, myCoverWidth, myCoverHeight, size, maxDepth, minArea)
		if err != nil {
			env.Printf("Error building cover: %s\n", err.Error())
			return nil, nil
		}
	}

	err = setEnvCover(env, cover)
	if err != nil {
		env.Printf("Error setting cover in project environment: %s\n", err.Error())
		return nil, nil
	}

	macro, img, err := findOrCreateMacro(env, cover, path, macroOutfile)
	if err != nil {
		env.Printf("Error building macro: %s\n", err.Error())
		coverService.Delete(cover)
		return cover, nil
	}

	err = setEnvMacro(env, macro)
	if err != nil {
		env.Printf("Error setting macro in project environment: %s\n", err.Error())
		return cover, nil
	}

	err = macroQuadBuildPartials(env, cover, macro, img, size, maxDepth, minArea)
	if err != nil {
		env.Printf("Error building quad partials: %s\n", err.Error())
		return cover, nil
	}

	if coverOutfile != "" {
		err = CoverDraw(env, cover.Id, coverOutfile)
		if err != nil {
			env.Printf("Error drawing cover: %s\n", err.Error())
			return cover, nil
		}
	}

	return cover, macro
}

func macroQuadBuildPartials(env environment.Environment, cover *model.Cover, macro *model.Macro, img *image.Image, size, maxDepth, minArea int) error {
	coverPartialService := env.MustCoverPartialService()
	quadDistService := env.MustQuadDistService()

	var err error

	total := 4 + 3*size
	count, err := coverPartialService.Count(cover)
	if err != nil {
		return err
	}
	current := int(count)
	remain := total - current

	msg := fmt.Sprintf("Building %d macro partial quads with size %d", remain, size)
	if maxDepth > 0 {
		msg += fmt.Sprintf(" and max depth %d", maxDepth)
	}
	if minArea > 0 {
		msg += fmt.Sprintf(" and min area %d", minArea)
	}
	msg += "..."
	env.Println(msg)

	bar := pb.StartNew(remain)

	for {
		var coverPartialQuadView *model.CoverPartialQuadView
		if current == 0 {
			// start with initial values
			coverPartialQuadView = &model.CoverPartialQuadView{
				CoverPartial: &model.CoverPartial{
					CoverId: cover.Id,
					X1:      0,
					Y1:      0,
					X2:      cover.Width,
					Y2:      cover.Height,
				},
				QuadDist: &model.QuadDist{
					Depth: 0,
					Area:  0,
					Dist:  0.0,
				},
			}
			current = 1 // fake the first
		} else {
			// get the worst quad dist
			coverPartialQuadView, err = quadDistService.GetWorst(macro, maxDepth, minArea)
			if err != nil {
				return err
			}

			if coverPartialQuadView == nil {
				return errors.New("Failed to find worst quad dist")
			}

			err = coverPartialService.Delete(coverPartialQuadView.CoverPartial)
			if err != nil {
				return err
			}
		}

		err = macroQuadSplit(env, macro, coverPartialQuadView, img)
		if err != nil {
			return err
		}
		current += 3
		bar.Set(current)

		if env.Cancel() {
			break
		}

		if current >= remain {
			break
		}
	}

	if env.Cancel() {
		return errors.New("Cancelled")
	}

	bar.Finish()
	return nil
}

func macroQuadCreateCover(env environment.Environment, width, height, size, maxDepth, minArea int) (*model.Cover, error) {
	aspectService := env.MustAspectService()
	coverService := env.MustCoverService()

	aspect, err := aspectService.FindOrCreate(width, height)
	if err != nil {
		return nil, err
	}

	cover := &model.Cover{
		AspectId: aspect.Id,
		Width:    width,
		Height:   height,
	}

	err = coverService.Insert(cover)
	if err != nil {
		return nil, err
	}

	return cover, nil
}

func macroQuadSplit(env environment.Environment, macro *model.Macro, coverPartialQuadView *model.CoverPartialQuadView, img *image.Image) error {
	coverPartials, err := macroQuadBuildCoverPartials(env, coverPartialQuadView)
	if err != nil {
		return err
	}

	macroPartials, err := macroQuadBuildMacroPartials(env, macro, coverPartials, img)
	if err != nil {
		return err
	}

	return macroQuadBuildQuadDist(env, coverPartials, macroPartials, coverPartialQuadView.QuadDist, img)
}

func macroQuadBuildCoverPartials(env environment.Environment, coverPartialQuadView *model.CoverPartialQuadView) ([]*model.CoverPartial, error) {
	aspectService := env.MustAspectService()
	coverPartialService := env.MustCoverPartialService()

	x1 := coverPartialQuadView.CoverPartial.X1
	y1 := coverPartialQuadView.CoverPartial.Y1
	x2 := coverPartialQuadView.CoverPartial.X2
	y2 := coverPartialQuadView.CoverPartial.Y2

	midX := ((x2 - x1) / 2) + x1
	midY := ((y2 - y1) / 2) + y1

	coverPartials := make([]*model.CoverPartial, 4)

	for i, pt := range [][]int{
		[]int{x1, y1, midX, midY},
		[]int{midX, y1, x2, midY},
		[]int{x1, midY, midX, y2},
		[]int{midX, midY, x2, y2},
	} {
		cp := &model.CoverPartial{
			CoverId: coverPartialQuadView.CoverPartial.CoverId,
			X1:      pt[0],
			Y1:      pt[1],
			X2:      pt[2],
			Y2:      pt[3],
		}

		aspect, err := aspectService.FindOrCreate(cp.Width(), cp.Height())
		if err != nil {
			return nil, err
		}

		cp.AspectId = aspect.Id

		err = coverPartialService.Insert(cp)
		if err != nil {
			return nil, err
		}

		coverPartials[i] = cp
	}

	return coverPartials, nil
}

func macroQuadBuildMacroPartials(env environment.Environment, macro *model.Macro, coverPartials []*model.CoverPartial, img *image.Image) ([]*model.MacroPartial, error) {
	macroPartialService := env.MustMacroPartialService()

	macroPartials := make([]*model.MacroPartial, 4)
	sem := make(chan bool, 4)

	for idx, coverPartial := range coverPartials {
		sem <- true
		go func(i int, cp *model.CoverPartial) {
			macroPartial := model.MacroPartial{
				MacroId:        macro.Id,
				CoverPartialId: cp.Id,
				AspectId:       cp.AspectId,
			}

			macroPartial.Pixels = util.GetImgPartialLab(img, cp)
			err := macroPartialService.Insert(&macroPartial)
			if err != nil {
				macroPartials[i] = nil
				<-sem
				return
			}

			macroPartials[i] = &macroPartial
			<-sem
		}(idx, coverPartial)
	}

	for i := 0; i < cap(sem); i++ {
		sem <- true
	}

	close(sem)

	for _, mp := range macroPartials {
		if mp == nil {
			return nil, errors.New("Failed to create macro partial")
		}
	}

	return macroPartials, nil
}

func macroQuadBuildQuadDist(env environment.Environment, coverPartials []*model.CoverPartial, macroPartials []*model.MacroPartial, parent *model.QuadDist, img *image.Image) error {
	quadDistService := env.MustQuadDistService()

	sem := make(chan bool, 4)
	errs := false

	for idx := 0; idx < 4; idx++ {
		sem <- true
		go func(i int) {
			quadDist := &model.QuadDist{
				MacroPartialId: macroPartials[i].Id,
				Depth:          parent.Depth + 1,
				Area:           coverPartials[i].Area(),
				Dist:           util.GetImgAvgDist(img, coverPartials[i]),
			}
			err := quadDistService.Insert(quadDist)
			if err != nil {
				errs = true
			}
			<-sem
		}(idx)
	}

	for i := 0; i < cap(sem); i++ {
		sem <- true
	}

	close(sem)

	if errs {
		return errors.New("Failed to create macro quad dist")
	}

	return nil
}

// for size, maxDepth, and minArea
// val == 0 is unrestricted
// val > 0 sets explicitly
// val == -1 (<0) calculates optimal
func macroQuadFixArgs(width, height, size, maxDepth, minArea int) (int, int, int) {
	var cSize, cMaxDepth, cMinArea int

	if size <= 0 {
		// set size to 2/5 root of total number of pixels
		area := width * height
		cSize = util.Round(math.Pow(float64(area), 0.4))
	} else {
		cSize = size
	}

	// aveDim is the average dimension of width and height
	aveDim := util.Round((float64(width) + float64(height)) / 2.0)

	if minArea <= 0 {
		// min size is the smallest length of a macro partial that we can tolerate
		// it is the bigger of size cut into 85 partials, and 35px
		minSize := util.Round(math.Max(float64(aveDim/85), float64(35)))
		cMinArea = minSize * minSize
	} else {
		cMinArea = minArea
	}

	if maxDepth <= 0 {
		// we want a max depth such that
		// size / 2^depth = minArea ^ (1/2)
		// solve for depth
		v1 := float64(aveDim * aveDim / cMinArea)
		cMaxDepth = util.Round(math.Log(v1) / (2.0 * math.Log(2.0)))
	} else {
		cMaxDepth = maxDepth
	}

	return cSize, cMaxDepth, cMinArea
}
