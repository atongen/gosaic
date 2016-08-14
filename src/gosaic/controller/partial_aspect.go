package controller

import (
	"errors"
	"gosaic/environment"
	"gosaic/model"
	"gosaic/service"
	"gosaic/util"
	"log"
)

func PartialAspect(env environment.Environment, dims ...int) {
	aspectService, err := env.AspectService()
	if err != nil {
		env.Fatalf("Error getting aspect service: %s\n", err.Error())
	}

	gidxPartialService, err := env.GidxPartialService()
	if err != nil {
		env.Fatalf("Error getting gidx partial service: %s\n", err.Error())
	}

	aspects, err := aspectsFromDims(aspectService, dims)
	if err != nil {
		env.Fatalf("Error getting aspects from dimesions: %s\n", err.Error())
	}

	err = createMissingGidxIndexes(env.Log(), gidxPartialService, aspects)
	if err != nil {
		env.Fatalf("Error creating missing indexes: %s\n", err.Error())
	}
}

func aspectsFromDims(aspectService service.AspectService, dims []int) ([]*model.Aspect, error) {
	if len(dims)%2 == 1 {
		return nil, errors.New("length of dimensions must be even, WxH,...")
	}

	aspects := make([]*model.Aspect, 0)
	aspectIds := make([]int64, 0)
	for i := 0; i < len(dims); i += 2 {
		w := dims[i]
		h := dims[i+1]

		aspect, err := aspectService.FindOrCreate(w, h)
		if err != nil {
			return nil, err
		}

		if !util.SliceContainsInt64(aspectIds, aspect.Id) {
			aspects = append(aspects, aspect)
			aspectIds = append(aspectIds, aspect.Id)
		}
	}

	return aspects, nil
}

func createMissingGidxIndexes(l *log.Logger, gidxPartialService service.GidxPartialService, aspects []*model.Aspect) error {
	for _, aspect := range aspects {
		for {
			gidxs, err := gidxPartialService.FindMissing(aspect, "gidx.id ASC", 1000, 0)
			if err != nil {
				return err
			}

			if len(gidxs) == 0 {
				break
			}

			l.Printf("Creating %d index partials for aspect %dx%d\n", len(gidxs), aspect.Columns, aspect.Rows)

			for _, gidx := range gidxs {
				_, err := gidxPartialService.Create(gidx, aspect)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
