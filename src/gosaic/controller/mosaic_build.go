package controller

import (
	"errors"
	"fmt"
	"gosaic/environment"
	"gosaic/model"
	"math"

	"gopkg.in/cheggaaa/pb.v1"
)

func MosaicBuild(env environment.Environment, name, fillType string, macroId int64, maxRepeats int) *model.Mosaic {
	gidxPartialService := env.MustGidxPartialService()
	macroService := env.MustMacroService()
	macroPartialService := env.MustMacroPartialService()
	mosaicService := env.MustMosaicService()

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
		err = createMosaicPartialsRandom(env, mosaic, maxRepeats)
	case "best":
		err = createMosaicPartialsBest(env, mosaic, maxRepeats)
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

func createMosaicPartialsRandom(env environment.Environment, mosaic *model.Mosaic, maxRepeats int) error {
	mosaicPartialService := env.MustMosaicPartialService()
	partialComparisonService := env.MustPartialComparisonService()

	numMissing, err := mosaicPartialService.CountMissing(mosaic)
	if err != nil {
		return err
	}

	if numMissing == 0 {
		return nil
	}

	env.Printf("Building %d mosaic partials...\n", numMissing)
	bar := pb.StartNew(int(numMissing))

	for {
		if env.Cancel() {
			return errors.New("Cancelled")
		}

		macroPartial, err := mosaicPartialService.GetRandomMissing(mosaic)
		if err != nil {
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
			return err
		}

		if gidxPartialId == int64(0) {
			return fmt.Errorf("Error: Invalid closest gidx partial found")
		}

		mosaicPartial := model.MosaicPartial{
			MosaicId:       mosaic.Id,
			MacroPartialId: macroPartial.Id,
			GidxPartialId:  gidxPartialId,
		}
		err = mosaicPartialService.Insert(&mosaicPartial)
		if err != nil {
			return err
		}
		bar.Increment()
	}

	bar.Finish()
	return nil
}

func createMosaicPartialsBest(env environment.Environment, mosaic *model.Mosaic, maxRepeats int) error {
	mosaicPartialService := env.MustMosaicPartialService()
	partialComparisonService := env.MustPartialComparisonService()

	numMissing, err := mosaicPartialService.CountMissing(mosaic)
	if err != nil {
		return err
	}

	if numMissing == 0 {
		return nil
	}

	env.Printf("Building %d mosaic partials...\n", numMissing)
	bar := pb.StartNew(int(numMissing))

	for {
		if env.Cancel() {
			return errors.New("Cancelled")
		}

		var partialComparison *model.PartialComparison
		if maxRepeats == 0 {
			partialComparison, err = partialComparisonService.GetBestAvailable(mosaic)
		} else {
			partialComparison, err = partialComparisonService.GetBestAvailableMax(mosaic, maxRepeats)
		}

		if err != nil {
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
			return err
		}
		bar.Increment()
	}

	bar.Finish()
	return nil
}
