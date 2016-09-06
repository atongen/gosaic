package controller

import (
	"errors"
	"fmt"
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

func MosaicBuild(env environment.Environment, name, fillType string, macroId int64, maxRepeats int) *model.Mosaic {
	gidxPartialService, err := env.GidxPartialService()
	if err != nil {
		env.Printf("Error getting index partial service: %s\n", err.Error())
		return nil
	}

	macroService, err := env.MacroService()
	if err != nil {
		env.Printf("Error getting macro service: %s\n", err.Error())
		return nil
	}

	macroPartialService, err := env.MacroPartialService()
	if err != nil {
		env.Printf("Error getting macro partial service: %s\n", err.Error())
		return nil
	}

	mosaicService, err := env.MosaicService()
	if err != nil {
		env.Printf("Error getting mosaic service: %s\n", err.Error())
		return nil
	}

	mosaicPartialService, err := env.MosaicPartialService()
	if err != nil {
		env.Printf("Error getting mosaic partial service: %s\n", err.Error())
		return nil
	}

	partialComparisonService, err := env.PartialComparisonService()
	if err != nil {
		env.Printf("Error getting partial comparison service: %s\n", err.Error())
		return nil
	}

	macro, err := macroService.Get(macroId)
	if err != nil {
		env.Printf("Error getting macro: %s\n", err.Error())
		return nil
	}

	if macro == nil {
		env.Printf("Macro %d not found\n", macroId)
		return nil
	}

	numMacroPartials, err := macroPartialService.Count(macro)
	if err != nil {
		env.Printf("Error counting macro partials: %s\n", err.Error())
		return nil
	}

	numGidxs, err := gidxPartialService.CountForMacro(macro)
	if err != nil {
		env.Printf("Error counting index images: %s\n", err.Error())
		return nil
	}

	// maxRepeats == 0 is unrestricted
	// maxRepeats > 0 sets explicitly
	// maxRepeats == -1 (<0) calculates minimum
	if maxRepeats < 0 {
		if numGidxs >= numMacroPartials {
			maxRepeats = 1
		} else {
			maxRepeats = int(math.Ceil(float64(numMacroPartials) / float64(numGidxs)))
		}
	} else if maxRepeats > 0 {
		if numGidxs*int64(maxRepeats) < numMacroPartials {
			env.Printf("Not enough index images (%d) to fill mosaic (%d) with max repeats set to %d", numGidxs, numMacroPartials, maxRepeats)
			return nil
		}
	}

	mosaic, err := mosaicService.GetOneBy("macro_id = ? AND name = ?", macroId, name)
	if err != nil {
		env.Printf("Error checking for existing mosaic: %s\n", err.Error())
		return nil
	}

	if mosaic == nil {
		mosaic = &model.Mosaic{
			Name:    name,
			MacroId: macro.Id,
		}
		err = mosaicService.Insert(mosaic)
		if err != nil {
			env.Printf("Error creating mosaic: %s\n", err.Error())
			return nil
		}
	}

	switch fillType {
	default:
		env.Printf("Invalid mosaic type: %s\n", fillType)
		return nil
	case "random":
		err = createMosaicPartialsRandom(env.Log(), mosaicPartialService, partialComparisonService, mosaic, maxRepeats)
	case "best":
		err = createMosaicPartialsBest(env.Log(), mosaicPartialService, partialComparisonService, mosaic, maxRepeats)
	}
	if err != nil {
		env.Printf("Error creating mosaic partials: %s\n", err.Error())
		return nil
	}

	mosaic.IsComplete = true
	_, err = mosaicService.Update(mosaic)
	if err != nil {
		env.Printf("Error marking mosaic complete: %s\n", err.Error())
		return nil
	}

	return mosaic
}

func createMosaicPartialsRandom(l *log.Logger, mosaicPartialService service.MosaicPartialService, partialComparisonService service.PartialComparisonService, mosaic *model.Mosaic, maxRepeats int) error {
	numMissing, err := mosaicPartialService.CountMissing(mosaic)
	if err != nil {
		return err
	}

	if numMissing == 0 {
		return nil
	}

	l.Printf("Building %d mosaic partials...\n", numMissing)
	bar := pb.StartNew(int(numMissing))

	cancel := false
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cancel = true
	}()

	for {
		if cancel {
			close(c)
			return errors.New("Cancelled")
		}

		macroPartial, err := mosaicPartialService.GetRandomMissing(mosaic)
		if err != nil {
			close(c)
			return err
		}
		if macroPartial == nil {
			break
		}

		var gidxPartialId int64
		if maxRepeats == 0 {
			gidxPartialId, err = partialComparisonService.GetClosest(macroPartial)
		} else {
			gidxPartialId, err = partialComparisonService.GetClosestMax(macroPartial, mosaic, maxRepeats)
		}
		if err != nil {
			close(c)
			return err
		}

		if gidxPartialId == int64(0) {
			close(c)
			return fmt.Errorf("Error: Invalid closest gidx partial found")
		}

		mosaicPartial := model.MosaicPartial{
			MosaicId:       mosaic.Id,
			MacroPartialId: macroPartial.Id,
			GidxPartialId:  gidxPartialId,
		}
		err = mosaicPartialService.Insert(&mosaicPartial)
		if err != nil {
			close(c)
			return err
		}
		bar.Increment()
	}

	close(c)
	bar.Finish()
	return nil
}

func createMosaicPartialsBest(l *log.Logger, mosaicPartialService service.MosaicPartialService, partialComparisonService service.PartialComparisonService, mosaic *model.Mosaic, maxRepeats int) error {
	numMissing, err := mosaicPartialService.CountMissing(mosaic)
	if err != nil {
		return err
	}

	if numMissing == 0 {
		return nil
	}

	l.Printf("Building %d mosaic partials...\n", numMissing)
	bar := pb.StartNew(int(numMissing))

	cancel := false
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cancel = true
	}()

	for {
		if cancel {
			close(c)
			return errors.New("Cancelled")
		}

		var partialComparison *model.PartialComparison
		if maxRepeats == 0 {
			partialComparison, err = partialComparisonService.GetBestAvailable(mosaic)
		} else {
			partialComparison, err = partialComparisonService.GetBestAvailableMax(mosaic, maxRepeats)
		}

		if err != nil {
			close(c)
			return err
		} else if partialComparison == nil {
			break
		}

		mosaicPartial := model.MosaicPartial{
			MosaicId:       mosaic.Id,
			MacroPartialId: partialComparison.MacroPartialId,
			GidxPartialId:  partialComparison.GidxPartialId,
		}
		err = mosaicPartialService.Insert(&mosaicPartial)
		if err != nil {
			close(c)
			return err
		}
		bar.Increment()
	}

	close(c)
	bar.Finish()
	return nil
}
