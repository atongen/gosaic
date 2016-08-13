package controller

import (
	"gosaic/environment"
	"gosaic/model"
	"gosaic/service"
	"log"
)

func Compare(env environment.Environment, macroId int64) {
	aspectService, err := env.AspectService()
	if err != nil {
		env.Fatalf("Error getting aspect service: %s\n", err.Error())
	}

	gidxPartialService, err := env.GidxPartialService()
	if err != nil {
		env.Fatalf("Error getting gidx partial service: %s\n", err.Error())
	}

	macroService, err := env.MacroService()
	if err != nil {
		env.Fatalf("Error getting macro service: %s\n", err.Error())
	}

	macroPartialService, err := env.MacroPartialService()
	if err != nil {
		env.Fatalf("Error getting macro partial service: %s\n", err.Error())
	}

	partialComparisonService, err := env.PartialComparisonService()
	if err != nil {
		env.Fatalf("Error getting partial comparison service: %s\n", err.Error())
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

	createMissingComparisons(env.Log(), partialComparisonService, macro)
}

func createMissingComparisons(l *log.Logger, partialComparisonService service.PartialComparisonService, macro *model.Macro) {
	batchSize := 100
	numTotal, err := partialComparisonService.CountMissing(macro)
	if err != nil {
		l.Fatalf("Error counting missing partial comparisons %s\n", err.Error())
	}
	created := int64(0)

	for {
		views, err := partialComparisonService.FindMissing(macro, batchSize)
		if err != nil {
			l.Fatalf("Error finding missing comparisons: %s\n", err.Error())
		}

		if len(views) == 0 {
			break
		}

		partialComparisons, err := buildPartialComparisons(views)
		if err != nil {
			l.Fatalf("Error building partial comparisons: %s\n", err.Error())
		}

		numCreated, err := partialComparisonService.BulkInsert(partialComparisons)
		if err != nil {
			l.Fatalf("Error inserting partial comparisons: %s\n", err.Error())
		}

		created += numCreated
		l.Printf("%d / %d partial comparisons created\n", created, numTotal)
	}
}

func buildPartialComparisons(macroGidxViews []*model.MacroGidxView) ([]*model.PartialComparison, error) {
	// TODO
	return nil, nil
}

func buildPartialComparison(macroGidxView *model.MacroGidxView) (*model.PartialComparison, error) {
	err := macroGidxView.MacroPartial.DecodeData()
	if err != nil {
		return nil, err
	}

	err = macroGidxView.GidxPartial.DecodeData()
	if err != nil {
		return nil, err
	}

	dist, err := model.PixelDist(macroGidxView.MacroPartial, macroGidxView.GidxPartial)
	if err != nil {
		return nil, err
	}

	return &model.PartialComparison{
		MacroPartialId: macroGidxView.MacroPartial.Id,
		GidxPartialId:  macroGidxView.GidxPartial.Id,
		Dist:           dist,
	}, nil

}
