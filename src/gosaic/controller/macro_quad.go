package controller

import (
	"errors"
	"gosaic/environment"
	"gosaic/model"
	"gosaic/service"
	"gosaic/util"
	"image"
	"log"
	"math"
	"os"
	"os/signal"
	"syscall"

	"gopkg.in/cheggaaa/pb.v1"
)

func MacroQuad(env environment.Environment,
	path string,
	coverWidth, coverHeight, num, maxDepth, minArea int,
	coverOutfile, macroOutfile string) (*model.Cover, *model.Macro) {

	aspectService, err := env.AspectService()
	if err != nil {
		env.Printf("Error getting aspect service: %s\n", err.Error())
		return nil, nil
	}

	coverService, err := env.CoverService()
	if err != nil {
		env.Printf("Error creating cover service: %s\n", err.Error())
		return nil, nil
	}

	coverPartialService, err := env.CoverPartialService()
	if err != nil {
		env.Printf("Error creating cover partial service: %s\n", err.Error())
		return nil, nil
	}

	macroService, err := env.MacroService()
	if err != nil {
		env.Printf("Error creating macro service: %s\n", err.Error())
		return nil, nil
	}

	macroPartialService, err := env.MacroPartialService()
	if err != nil {
		env.Printf("Error creating macro partial service: %s\n", err.Error())
		return nil, nil
	}

	quadDistService, err := env.QuadDistService()
	if err != nil {
		env.Printf("Error creating quad dist service: %s\n", err.Error())
		return nil, nil
	}

	myCoverWidth, myCoverHeight, err := calculateDimensions(aspectService, path, coverWidth, coverHeight)
	if err != nil {
		env.Printf("Error getting cover dimensions: %s\n", err.Error())
		return nil, nil
	}

	num, maxDepth, minArea = macroQuadFixArgs(myCoverWidth, myCoverHeight, num, maxDepth, minArea)

	cover, err := macroQuadCreateCover(coverService, aspectService, myCoverWidth, myCoverHeight, num, maxDepth, minArea)
	if err != nil {
		env.Printf("Error building cover: %s\n", err.Error())
		return nil, nil
	}

	macro, img, err := findOrCreateMacro(macroService, aspectService, cover, path, macroOutfile)
	if err != nil {
		env.Printf("Error building macro: %s\n", err.Error())
		coverService.Delete(cover)
		return nil, nil
	}

	err = macroQuadBuildPartials(env.Log(), aspectService, coverPartialService, macroPartialService, quadDistService, cover, macro, img, num, maxDepth, minArea)
	if err != nil {
		env.Printf("Error building quad partials: %s\n", err.Error())
		coverService.Delete(cover)
		return nil, nil
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

func macroQuadBuildPartials(l *log.Logger, aspectService service.AspectService, coverPartialService service.CoverPartialService, macroPartialService service.MacroPartialService, quadDistService service.QuadDistService, cover *model.Cover, macro *model.Macro, img *image.Image, num, maxDepth, minArea int) error {
	var err error

	// start with initial values
	coverPartialQuadView := &model.CoverPartialQuadView{
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

	l.Printf("Building %d macro partial quads...\n", num)

	bar := pb.StartNew(num)

	cancel := false
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cancel = true
	}()

	for i := 0; i < num; i++ {
		if cancel {
			break
		}

		err = macroQuadSplit(aspectService, coverPartialService, macroPartialService, quadDistService, macro, coverPartialQuadView, img)
		if err != nil {
			return err
		}

		coverPartialQuadView, err = quadDistService.GetWorst(macro, maxDepth, minArea)
		if err != nil {
			return err
		}

		if coverPartialQuadView == nil {
			// we are done
			break
		}

		if coverPartialQuadView.CoverPartial.Id != int64(0) {
			err = coverPartialService.Delete(coverPartialQuadView.CoverPartial)
			if err != nil {
				return err
			}
		}

		bar.Increment()
	}

	if cancel {
		return errors.New("Cancelled")
	}

	bar.Finish()
	return nil
}

func macroQuadCreateCover(coverService service.CoverService, aspectService service.AspectService, width, height, num, maxDepth, minArea int) (*model.Cover, error) {
	aspect, err := aspectService.FindOrCreate(width, height)
	if err != nil {
		return nil, err
	}

	coverName := model.CoverNameQuad(aspect.Id, width, height, num, maxDepth, minArea)
	cover := model.Cover{
		Name:     coverName,
		AspectId: aspect.Id,
		Width:    width,
		Height:   height,
	}

	err = coverService.Insert(&cover)
	if err != nil {
		return nil, err
	}

	return &cover, nil
}

func macroQuadSplit(aspectService service.AspectService, coverPartialService service.CoverPartialService, macroPartialService service.MacroPartialService, quadDistService service.QuadDistService, macro *model.Macro, coverPartialQuadView *model.CoverPartialQuadView, img *image.Image) error {
	coverPartials, err := macroQuadBuildCoverPartials(coverPartialService, aspectService, coverPartialQuadView)
	if err != nil {
		return err
	}

	macroPartials, err := macroQuadBuildMacroPartials(macroPartialService, macro, coverPartials, img)
	if err != nil {
		return err
	}

	return macroQuadBuildQuadDist(quadDistService, coverPartials, macroPartials, coverPartialQuadView.QuadDist, img)
}

func macroQuadBuildCoverPartials(coverPartialService service.CoverPartialService, aspectService service.AspectService, coverPartialQuadView *model.CoverPartialQuadView) ([]*model.CoverPartial, error) {
	x1 := coverPartialQuadView.CoverPartial.X1
	y1 := coverPartialQuadView.CoverPartial.Y1
	x2 := coverPartialQuadView.CoverPartial.X2
	y2 := coverPartialQuadView.CoverPartial.Y2

	midX := ((x2 - x1) / 2) + x1
	midY := ((y2 - y1) / 2) + y1

	coverPartials := make([]*model.CoverPartial, 4)

	for i, pt := range [][]int{
		[]int{x1, y1, midX, midY},
		[]int{midX + 1, y1, x2, midY},
		[]int{x1, midY + 1, midX, y2},
		[]int{midX + 1, midY + 1, x2, y2},
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

func macroQuadBuildMacroPartials(macroPartialService service.MacroPartialService, macro *model.Macro, coverPartials []*model.CoverPartial, img *image.Image) ([]*model.MacroPartial, error) {
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

func macroQuadBuildQuadDist(quadDistService service.QuadDistService, coverPartials []*model.CoverPartial, macroPartials []*model.MacroPartial, parent *model.QuadDist, img *image.Image) error {
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

func macroQuadFixArgs(width, height, num, maxDepth, minArea int) (int, int, int) {
	var size, cNum, cMaxDepth, cMinArea int

	if num > 0 {
		cNum = num
	} else {
		// arbitrarily choose n iterations if none provided
		cNum = 1024
	}

	// size is the smaller dimension of width and height
	if width < height {
		size = width
	} else {
		size = height
	}

	if minArea > 0 {
		cMinArea = minArea
	} else {
		// min size is the smallest length of a macro partial that we can tolerate
		// it is the bigger of size cut into 100 partials, and 25px
		minSize := util.Round(math.Max(float64(size/100), float64(25)))
		cMinArea = minSize * minSize
	}

	if maxDepth > 0 {
		cMaxDepth = maxDepth
	} else {
		// we want a max depth such that
		// size / 2^depth = minArea ^ (1/2)
		// solve for depth
		v1 := float64(size * size / cMinArea)
		cMaxDepth = util.Round(math.Log(v1) / (2.0 * math.Log(2.0)))
	}

	return cNum, cMaxDepth, cMinArea
}
