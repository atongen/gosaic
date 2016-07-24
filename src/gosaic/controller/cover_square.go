package controller

import (
	"fmt"
	"gosaic/environment"
	"gosaic/model"
	"gosaic/service"
	"math"
)

func CoverSquare(env environment.Environment, name string, width, height, num int) {
	coverService, err := env.CoverService()
	if err != nil {
		fmt.Printf("Error creating cover service: %s\n", err.Error())
		return
	}

	coverPartialService, err := env.CoverPartialService()
	if err != nil {
		fmt.Printf("Error creating cover partial service: %s\n", err.Error())
		return
	}

	aspectService, err := env.AspectService()
	if err != nil {
		fmt.Printf("Error getting aspect service: %s\n", err.Error())
		return
	}

	cover, err := createSquareCover(coverService, name, width, height)
	if err != nil {
		fmt.Printf("Error creating square cover: %s\n", err.Error())
		return
	}

	aspect, err := aspectService.FindOrCreate(1, 1)
	if err != nil {
		fmt.Printf("Error getting square aspect: %s\n", err.Error())
		return
	}

	numPartials, err := addSquareCoverPartials(coverPartialService, cover, aspect, num)
	if err != nil {
		fmt.Printf("Error adding square cover partials: %s\n", err.Error())
		return
	}

	fmt.Printf("Created cover %s with %d partials\n", cover.Name, numPartials)
}

func createSquareCover(coverService service.CoverService, name string, width, height int) (*model.Cover, error) {
	var cover model.Cover = model.Cover{
		Type:   "square",
		Name:   name,
		Width:  uint(width),
		Height: uint(height),
	}

	err := coverService.Insert(&cover)
	if err != nil {
		return nil, err
	}

	return &cover, nil
}

func getSquareCoverDims(width, height, num int) (size, columns, rows int) {
	fw := float64(width)
	fh := float64(height)
	fn := float64(num)

	var fs, fc, fr float64

	if fw < fh {
		fs = math.Ceil(fw / fn)
		fc = fn
		fr = math.Ceil(fh / fs)
	} else {
		fs = math.Ceil(fh / fn)
		fr = fn
		fc = math.Ceil(fw / fs)
	}

	size = int(fs)
	columns = int(fc)
	rows = int(fr)

	return
}

func addSquareCoverPartials(coverPartialService service.CoverPartialService, cover *model.Cover, aspect *model.Aspect, num int) (int, error) {
	size, columns, rows := getSquareCoverDims(int(cover.Width), int(cover.Height), num)

	xOffset := int(math.Floor(float64(int(cover.Width)-size*columns) / float64(2.0)))
	yOffset := int(math.Floor(float64(int(cover.Height)-size*rows) / float64(2.0)))

	created := 0
	total := columns * rows
	fmt.Printf("size: %d, columns: %d, rows: %d\n", size, columns, rows)
	fmt.Printf("xOffset: %d, yOffset: %d\n", xOffset, yOffset)

	for i := 0; i < columns; i++ {
		for j := 0; j < rows; j++ {
			x1 := i*size + xOffset
			y1 := j*size + yOffset
			x2 := (i+1)*size + xOffset - 1
			y2 := (j+1)*size + yOffset - 1

			var coverPartial model.CoverPartial = model.CoverPartial{
				CoverId:  cover.Id,
				AspectId: aspect.Id,
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
			fmt.Printf("%d/%d - x1: %d, y1: %d, x2: %d, y2: %d\n", created, total, x1, y1, x2, y2)
		}
	}

	return created, nil
}
