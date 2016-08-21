package controller

import (
	"gosaic/environment"
	"gosaic/model"
	"gosaic/service"
	"log"

	"gopkg.in/cheggaaa/pb.v1"
)

func PartialAspect(env environment.Environment, macroId int64) {
	aspectService, err := env.AspectService()
	if err != nil {
		env.Fatalf("Error getting aspect service: %s\n", err.Error())
	}

	macroService, err := env.MacroService()
	if err != nil {
		env.Fatalf("Error getting macro service: %s\n", err.Error())
	}

	macroPartialService, err := env.MacroPartialService()
	if err != nil {
		env.Fatalf("Error getting macro partial service: %s\n", err.Error())
	}

	gidxPartialService, err := env.GidxPartialService()
	if err != nil {
		env.Fatalf("Error getting gidx partial service: %s\n", err.Error())
	}

	macro, err := macroService.Get(macroId)
	if err != nil {
		env.Fatalf("Error getting macro %d: %s\n", macroId, err.Error())
	}

	aspectIds, err := macroPartialService.AspectIds(macro.Id)
	if err != nil {
		env.Fatalf("Error getting aspect ids: %s\n", err.Error())
	}

	if len(aspectIds) == 0 {
		env.Fatalln("No aspects found for this macro's partials")
	}

	aspects, err := aspectService.FindIn(aspectIds)
	if err != nil {
		env.Fatalf("Error getting aspects: %s\n", err.Error())
	}

	err = createMissingGidxIndexes(env.Log(), gidxPartialService, aspects)
	if err != nil {
		env.Fatalf("Error creating gidx partial aspects: %s\n", err.Error())
	}
}

func createMissingGidxIndexes(l *log.Logger, gidxPartialService service.GidxPartialService, aspects []*model.Aspect) error {
	count, err := gidxPartialService.CountMissing(aspects)
	if err != nil {
		return err
	}

	if count == 0 {
		return nil
	}

	l.Printf("Creating %d aspect partials for indexed images...\n", count)
	bar := pb.StartNew(int(count))

	for _, aspect := range aspects {
		for {
			gidxs, err := gidxPartialService.FindMissing(aspect, "gidx.id ASC", 100, 0)
			if err != nil {
				return err
			}

			if len(gidxs) == 0 {
				break
			}

			for _, gidx := range gidxs {
				_, err := gidxPartialService.Create(gidx, aspect)
				if err != nil {
					return err
				}
				bar.Increment()
			}
		}
	}

	bar.Finish()

	return nil
}
