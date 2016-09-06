package controller

import (
	"errors"
	"gosaic/environment"
	"gosaic/model"
	"gosaic/service"
	"log"
	"math"
	"os"
	"os/signal"
	"syscall"

	"gopkg.in/cheggaaa/pb.v1"
)

func CoverAspect(env environment.Environment, coverWidth, coverHeight, partialWidth, partialHeight, num int) *model.Cover {
	coverService, err := env.CoverService()
	if err != nil {
		env.Printf("Error creating cover service: %s\n", err.Error())
		return nil
	}

	coverPartialService, err := env.CoverPartialService()
	if err != nil {
		env.Printf("Error creating cover partial service: %s\n", err.Error())
		return nil
	}

	aspectService, err := env.AspectService()
	if err != nil {
		env.Printf("Error getting aspect service: %s\n", err.Error())
		return nil
	}

	coverAspect, err := aspectService.FindOrCreate(coverWidth, coverHeight)
	if err != nil {
		env.Printf("Error getting cover aspect: %s\n", err.Error())
		return nil
	}

	coverPartialAspect, err := aspectService.FindOrCreate(partialWidth, partialHeight)
	if err != nil {
		env.Printf("Error getting cover partial aspect: %s\n", err.Error())
		return nil
	}

	var cover *model.Cover
	coverName := model.CoverNameAspect(coverAspect.Id, coverWidth, coverHeight, num)
	cover, err = coverService.GetOneBy("name = ?", coverName)
	if err != nil {
		env.Printf("Error finding cover: %s\n", err.Error())
		return nil
	}
	// Existing cover is found, use it
	if cover != nil {
		return cover
	}

	cover = &model.Cover{
		Name:     coverName,
		AspectId: coverAspect.Id,
		Width:    coverWidth,
		Height:   coverHeight,
	}
	err = coverService.Insert(cover)
	if err != nil {
		env.Printf("Error creating cover: %s\n", err.Error())
		return nil
	}

	err = addCoverAspectPartials(env.Log(), coverPartialService, cover, coverPartialAspect, num)
	if err != nil {
		env.Printf("Error adding cover partials: %s\n", err.Error())
		// attempt to delete cover
		// this will fail if there is already a macro referencing it
		// which is fine
		coverService.Delete(cover)
		return nil
	}

	return cover
}

// getCoverAspectDims takes a cover width and height,
// and an expected width and height for the aspect of the partial rectangles
// needed and returns the width and height of the partial rectangles,
// as well as the number of columns and rows needed for the cover.
func getCoverAspectDims(coverWidth, coverHeight, partialAspectWidth, partialAspectHeight, num int) (width, height, columns, rows int) {
	cw := float64(coverWidth)
	ch := float64(coverHeight)
	aw := float64(partialAspectWidth)
	ah := float64(partialAspectHeight)
	n := float64(num)

	var fw, fh, pc, pr float64

	if cw < ch {
		fw = math.Ceil(cw / n)
		pc = n
		fh = math.Ceil(fw * ah / aw)
		pr = math.Ceil(ch / fh)
	} else {
		fh = math.Ceil(ch / n)
		pr = n
		fw = math.Ceil(fh * aw / ah)
		pc = math.Ceil(cw / fw)
	}

	width = int(fw)
	height = int(fh)
	columns = int(pc)
	rows = int(pr)

	return
}

func addCoverAspectPartials(l *log.Logger, coverPartialService service.CoverPartialService, cover *model.Cover, coverPartialAspect *model.Aspect, num int) error {
	width, height, columns, rows := getCoverAspectDims(cover.Width, cover.Height, coverPartialAspect.Columns, coverPartialAspect.Rows, num)

	xOffset := int(math.Floor(float64(cover.Width-width*columns) / float64(2.0)))
	yOffset := int(math.Floor(float64(cover.Height-height*rows) / float64(2.0)))

	count := columns * rows
	l.Printf("Building %d cover partials...\n", count)

	bar := pb.StartNew(count)

	cancel := false
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cancel = true
	}()

	for i := 0; i < columns; i++ {
		var coverPartials []*model.CoverPartial = make([]*model.CoverPartial, rows)
		for j := 0; j < rows; j++ {
			if cancel {
				close(c)
				return errors.New("Cancelled")
			}

			x1 := i*width + xOffset
			y1 := j*height + yOffset
			//x2 := (i+1)*width + xOffset - 1
			//y2 := (j+1)*height + yOffset - 1
			x2 := (i+1)*width + xOffset
			y2 := (j+1)*height + yOffset

			var coverPartial model.CoverPartial = model.CoverPartial{
				CoverId:  cover.Id,
				AspectId: coverPartialAspect.Id,
				X1:       x1,
				Y1:       y1,
				X2:       x2,
				Y2:       y2,
			}
			coverPartials[j] = &coverPartial
		}
		num, err := coverPartialService.BulkInsert(coverPartials)
		if err != nil {
			close(c)
			return err
		}
		bar.Add(int(num))
	}

	close(c)
	bar.Finish()
	return nil
}
