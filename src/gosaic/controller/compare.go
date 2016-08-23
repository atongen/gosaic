package controller

import (
	"errors"
	"gosaic/environment"
	"gosaic/model"
	"gosaic/service"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"gopkg.in/cheggaaa/pb.v1"
)

func Compare(env environment.Environment, macroId int64) error {
	macroService, err := env.MacroService()
	if err != nil {
		env.Printf("Error getting macro service: %s\n", err.Error())
		return err
	}

	partialComparisonService, err := env.PartialComparisonService()
	if err != nil {
		env.Printf("Error getting partial comparison service: %s\n", err.Error())
		return err
	}

	macro, err := macroService.Get(macroId)
	if err != nil {
		env.Printf("Error getting macro: %s\n", err.Error())
		return err
	}

	err = createMissingComparisons(env.Log(), partialComparisonService, macro)
	if err != nil {
		env.Printf("Error creating comparisons: %s\n", err.Error())
		return err
	}

	return nil
}

func createMissingComparisons(l *log.Logger, partialComparisonService service.PartialComparisonService, macro *model.Macro) error {
	batchSize := 500
	numTotal, err := partialComparisonService.CountMissing(macro)
	if err != nil {
		return err
	}

	if numTotal == 0 {
		return nil
	}

	l.Printf("Building %d partial image comparisons...\n", numTotal)
	bar := pb.StartNew(int(numTotal))

	cancel := false
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cancel = true
	}()

	for {
		if cancel {
			return errors.New("Cancelled")
		}

		views, err := partialComparisonService.FindMissing(macro, batchSize)
		if err != nil {
			return err
		}

		if len(views) == 0 {
			break
		}

		partialComparisons := buildPartialComparisons(l, views)

		if len(partialComparisons) > 0 {
			numCreated, err := partialComparisonService.BulkInsert(partialComparisons)
			if err != nil {
				return err
			}

			bar.Add(int(numCreated))
		}
	}
	bar.Finish()
	return nil
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
