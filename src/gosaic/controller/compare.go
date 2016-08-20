package controller

import (
	"gosaic/environment"
	"gosaic/model"
	"gosaic/service"
	"log"
	"sync"

	"gopkg.in/cheggaaa/pb.v1"
)

func Compare(env environment.Environment, macroId int64) {
	macroService, err := env.MacroService()
	if err != nil {
		env.Fatalf("Error getting macro service: %s\n", err.Error())
	}

	partialComparisonService, err := env.PartialComparisonService()
	if err != nil {
		env.Fatalf("Error getting partial comparison service: %s\n", err.Error())
	}

	macro, err := macroService.Get(macroId)
	if err != nil {
		env.Fatalf("Error getting macro: %s\n", err.Error())
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

	l.Printf("Creating %d partial image comparisons...\n", numTotal)
	bar := pb.StartNew(int(numTotal))

	for {
		views, err := partialComparisonService.FindMissing(macro, batchSize)
		if err != nil {
			l.Fatalf("Error finding missing comparisons: %s\n", err.Error())
		}

		if len(views) == 0 {
			break
		}

		partialComparisons := buildPartialComparisons(l, views)

		if len(partialComparisons) > 0 {
			numCreated, err := partialComparisonService.BulkInsert(partialComparisons)
			if err != nil {
				l.Fatalf("Error inserting partial comparisons: %s\n", err.Error())
			}

			bar.Add(int(numCreated))
		}
	}
	bar.Finish()
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
