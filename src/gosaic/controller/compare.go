package controller

import (
	"gosaic/environment"
	"gosaic/model"
	"gosaic/service"
	"log"
	"sync"
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

	if numTotal == 0 {
		l.Printf("No missing comparisons for macro %d\n", macro.Id)
		return
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

		partialComparisons := buildPartialComparisons(l, views)

		numErrors := len(views) - len(partialComparisons)
		if numErrors > 0 {
			l.Printf("Failed to create %d of %d comparisons\n", numErrors, len(views))
		}

		if len(partialComparisons) > 0 {
			numCreated, err := partialComparisonService.BulkInsert(partialComparisons)
			if err != nil {
				l.Fatalf("Error inserting partial comparisons: %s\n", err.Error())
			}

			created += numCreated
			l.Printf("%d / %d partial comparisons created\n", created, numTotal)
		}
	}
}

func buildPartialComparisons(l *log.Logger, macroGidxViews []*model.MacroGidxView) []*model.PartialComparison {
	var wg sync.WaitGroup
	wg.Add(len(macroGidxViews))

	partialComparisons := make([]*model.PartialComparison, len(macroGidxViews))
	add := make(chan *model.PartialComparison)
	errs := make(chan error)

	go func(myLog *log.Logger, pcs []*model.PartialComparison, addCh chan *model.PartialComparison, errsCh chan error) {
		idx := 0
		for i := 0; i < len(pcs); i++ {
			select {
			case pc := <-addCh:
				pcs[idx] = pc
				idx++
			case err := <-errsCh:
				l.Printf("Error building partial comparison: %s\n", err.Error())
			}
			wg.Done()
		}
	}(l, partialComparisons, add, errs)

	for _, macroGidxView := range macroGidxViews {
		go func(mgv *model.MacroGidxView, addCh chan *model.PartialComparison, errsCh chan error) {
			pc, err := mgv.PartialComparison()
			if err != nil {
				errsCh <- err
				return
			}
			addCh <- pc
		}(macroGidxView, add, errs)
	}

	wg.Wait()
	close(add)
	close(errs)
	return partialComparisons
}
