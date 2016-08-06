package controller

import (
	"gosaic/environment"
	"gosaic/model"
	"gosaic/service"
	"math"
)

func CoverAspect(env environment.Environment, name string, coverWidth, coverHeight, partialWidth, partialHeight, num int) {
	coverService, err := env.CoverService()
	if err != nil {
		env.Printf("Error creating cover service: %s\n", err.Error())
		return
	}

	coverPartialService, err := env.CoverPartialService()
	if err != nil {
		env.Printf("Error creating cover partial service: %s\n", err.Error())
		return
	}

	aspectService, err := env.AspectService()
	if err != nil {
		env.Printf("Error getting aspect service: %s\n", err.Error())
		return
	}

	coverAspect, err := aspectService.FindOrCreate(coverWidth, coverHeight)
	if err != nil {
		env.Printf("Error getting cover aspect: %s\n", err.Error())
		return
	}

	coverPartialAspect, err := aspectService.FindOrCreate(partialWidth, partialHeight)
	if err != nil {
		env.Printf("Error getting cover partial aspect: %s\n", err.Error())
		return
	}

	cover, err := createCoverAspect(coverService, coverAspect, name, coverWidth, coverHeight)
	if err != nil {
		env.Printf("Error creating aspect cover: %s\n", err.Error())
		return
	}

	numPartials, err := addCoverAspectPartials(coverPartialService, cover, coverPartialAspect, num)
	if err != nil {
		env.Printf("Error adding aspect cover partials: %s\n", err.Error())
		return
	}

	env.Printf("Created cover %s with %d partials\n", cover.Name, numPartials)
}

func createCoverAspect(coverService service.CoverService, aspect *model.Aspect, name string, width, height int) (*model.Cover, error) {
	var cover model.Cover = model.Cover{
		Type:     "aspect",
		AspectId: aspect.Id,
		Name:     name,
		Width:    uint(width),
		Height:   uint(height),
	}

	err := coverService.Insert(&cover)
	if err != nil {
		return nil, err
	}

	return &cover, nil
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

func addCoverAspectPartials(coverPartialService service.CoverPartialService, cover *model.Cover, coverPartialAspect *model.Aspect, num int) (int, error) {
	width, height, columns, rows := getCoverAspectDims(int(cover.Width), int(cover.Height), coverPartialAspect.Columns, coverPartialAspect.Rows, num)

	xOffset := int(math.Floor(float64(int(cover.Width)-width*columns) / float64(2.0)))
	yOffset := int(math.Floor(float64(int(cover.Height)-height*rows) / float64(2.0)))

	created := 0

	for i := 0; i < columns; i++ {
		for j := 0; j < rows; j++ {
			x1 := i*width + xOffset
			y1 := j*height + yOffset
			x2 := (i+1)*width + xOffset - 1
			y2 := (j+1)*height + yOffset - 1

			var coverPartial model.CoverPartial = model.CoverPartial{
				CoverId:  cover.Id,
				AspectId: coverPartialAspect.Id,
				X1:       int64(x1),
				Y1:       int64(y1),
				X2:       int64(x2),
				Y2:       int64(y2),
			}
			err := coverPartialService.Insert(&coverPartial)
			if err != nil {
				return created, err
			}
			created += 1
		}
	}

	return created, nil
}
