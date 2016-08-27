package controller

import (
	"errors"
	"gosaic/environment"
	"gosaic/model"
	"gosaic/service"
	"gosaic/util"
	"image"
	"log"
	"os"
	"os/signal"
	"syscall"

	"gopkg.in/cheggaaa/pb.v1"

	"github.com/disintegration/imaging"
)

func MacroQuad(env environment.Environment, path string, coverWidth, coverHeight, num int, outfile string) (*model.Cover, *model.Macro) {
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

	cover, err := macroQuadCreateCover(coverService, aspectService, myCoverWidth, myCoverHeight, num)
	if err != nil {
		env.Printf("Error building cover: %s\n", err.Error())
		return nil, nil
	}

	md5sum, err := util.Md5sum(path)
	if err != nil {
		env.Printf("Error getting macro md5sum: %s\n", err.Error())
		coverService.Delete(cover)
		return nil, nil
	}

	img, err := util.OpenImage(path)
	if err != nil {
		env.Printf("Failed to open image: %s\n", err.Error())
		coverService.Delete(cover)
		return nil, nil
	}

	orientation, err := util.GetOrientation(path)
	if err != nil {
		env.Printf("Failed to get image orientation: %s\n", err.Error())
		coverService.Delete(cover)
		return nil, nil
	}

	err = util.FixOrientation(img, orientation)
	if err != nil {
		env.Printf("Failed to fix image orientation: %s\n", err.Error())
		coverService.Delete(cover)
		return nil, nil
	}
	bounds := (*img).Bounds()

	var imgCov image.Image
	imgCov = imaging.Fill(*img, cover.Width, cover.Height, imaging.Center, imaging.Lanczos)

	if outfile != "" {
		err = imaging.Save(imgCov, outfile)
		if err != nil {
			env.Printf("Error saving file: %s\n", err.Error())
			coverService.Delete(cover)
			return nil, nil
		}
	}

	macroAspect, err := aspectService.FindOrCreate(bounds.Max.X, bounds.Max.Y)
	if err != nil {
		env.Printf("Failed to get macro aspect: %s\n", err.Error())
		coverService.Delete(cover)
		return nil, nil
	}

	macro := &model.Macro{
		AspectId:    macroAspect.Id,
		CoverId:     cover.Id,
		Path:        path,
		Md5sum:      md5sum,
		Width:       bounds.Max.X,
		Height:      bounds.Max.Y,
		Orientation: orientation,
	}
	err = macroService.Insert(macro)
	if err != nil {
		env.Printf("Error creating macro: %s\n", err.Error())
		coverService.Delete(cover)
		return nil, nil
	}

	err = macroQuadBuildPartials(env.Log(), aspectService, coverPartialService, macroPartialService, quadDistService, cover, macro, &imgCov, num)
	if err != nil {
		env.Printf("Error building quad partials: %s\n", err.Error())
		coverService.Delete(cover)
		return nil, nil
	}

	return cover, macro
}

func macroQuadBuildPartials(l *log.Logger, aspectService service.AspectService, coverPartialService service.CoverPartialService, macroPartialService service.MacroPartialService, quadDistService service.QuadDistService, cover *model.Cover, macro *model.Macro, img *image.Image, num int) error {
	var (
		err          error
		coverPartial *model.CoverPartial
	)

	coverPartial = &model.CoverPartial{
		X1: 0,
		Y1: 0,
		X2: cover.Width,
		Y2: cover.Height,
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

		err = macroQuadBuildFour(aspectService, coverPartialService, macroPartialService, quadDistService, cover, macro, img, coverPartial.X1, coverPartial.Y1, coverPartial.X2, coverPartial.Y2)
		if err != nil {
			return err
		}

		coverPartial, err = quadDistService.GetWorst(macro)
		if err != nil {
			return err
		}

		if coverPartial.Id != int64(0) {
			err = coverPartialService.Delete(coverPartial)
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

func macroQuadCreateCover(coverService service.CoverService, aspectService service.AspectService, width, height, num int) (*model.Cover, error) {
	aspect, err := aspectService.FindOrCreate(width, height)
	if err != nil {
		return nil, err
	}

	coverName := model.CoverNameQuad(aspect.Id, width, height, num)
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

func macroQuadBuildFour(aspectService service.AspectService, coverPartialService service.CoverPartialService, macroPartialService service.MacroPartialService, quadDistService service.QuadDistService, cover *model.Cover, macro *model.Macro, img *image.Image, x1, y1, x2, y2 int) error {
	coverPartials, err := macroQuadBuildCoverPartials(coverPartialService, aspectService, cover, x1, y1, x2, y2)
	if err != nil {
		return err
	}

	macroPartials, err := macroQuadBuildMacroPartials(macroPartialService, macro, coverPartials, img)
	if err != nil {
		return err
	}

	return macroQuadBuildQuadDist(quadDistService, coverPartials, macroPartials, img)
}

func macroQuadBuildCoverPartials(coverPartialService service.CoverPartialService, aspectService service.AspectService, cover *model.Cover, x1, y1, x2, y2 int) ([]*model.CoverPartial, error) {
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
			CoverId: cover.Id,
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

func macroQuadBuildQuadDist(quadDistService service.QuadDistService, coverPartials []*model.CoverPartial, macroPartials []*model.MacroPartial, img *image.Image) error {
	sem := make(chan bool, 4)
	errs := false

	for idx := 0; idx < 4; idx++ {
		sem <- true
		go func(i int) {
			quadDist := &model.QuadDist{
				MacroPartialId: macroPartials[i].Id,
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
