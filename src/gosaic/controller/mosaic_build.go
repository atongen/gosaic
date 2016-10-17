package controller

import (
	"errors"
	"fmt"
	"gosaic/environment"
	"gosaic/model"
	"math"

	"gopkg.in/cheggaaa/pb.v1"
)

func MosaicBuild(env environment.Environment, fillType string, macroId int64, maxRepeats int, destructive bool) *model.Mosaic {
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

	mosaic, err := envMosaic(env)
	if err != nil {
		env.Printf("Error getting mosaic from project environment: %s\n", err.Error())
		return nil
	}

	if mosaic == nil {
		mosaic = &model.Mosaic{
			MacroId: macro.Id,
		}
		err = mosaicService.Insert(mosaic)
		if err != nil {
			env.Printf("Error creating mosaic: %s\n", err.Error())
			return nil
		}
	}

	err = setEnvMosaic(env, mosaic)
	if err != nil {
		env.Printf("Error setting mosaic in project environment: %s\n", err.Error())
		return nil
	}

	err = doMosaicBuild(env, mosaic, fillType, maxRepeats, destructive)
	if err != nil {
		env.Printf("Error building mosaic: %s\n", err.Error())
		return nil
	}

	return mosaic
}

func doMosaicBuild(env environment.Environment, mosaic *model.Mosaic, fillType string, maxRepeats int, destructive bool) error {
	var err error
	switch fillType {
	default:
		env.Printf("Invalid mosaic type: %s\n", fillType)
		return nil
	case "random":
		err = createMosaicPartialsRandom(env, mosaic, maxRepeats, destructive)
	case "best":
		err = createMosaicPartialsBest(env, mosaic, maxRepeats, destructive)
	}
	if err != nil {
		return err
	}

	return nil
}

func createMosaicPartialsRandom(env environment.Environment, mosaic *model.Mosaic, maxRepeats int, destructive bool) error {
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
		if maxRepeats == 0 || destructive {
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

		if destructive && maxRepeats > 0 {
			err = mosaicBuildDeleteDuplicates(env, mosaic, maxRepeats)
			if err != nil {
				return err
			}
		}

		bar.Increment()
	}

	bar.Finish()
	return nil
}

func createMosaicPartialsBest(env environment.Environment, mosaic *model.Mosaic, maxRepeats int, destructive bool) error {
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
		if maxRepeats == 0 || destructive {
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

		if destructive && maxRepeats > 0 {
			err = mosaicBuildDeleteDuplicates(env, mosaic, maxRepeats)
			if err != nil {
				return err
			}
		}

		bar.Increment()
	}

	bar.Finish()
	return nil
}

func mosaicBuildDeleteDuplicates(env environment.Environment, mosaic *model.Mosaic, maxRepeats int) error {
	mosaicPartialService := env.MustMosaicPartialService()

	macroPartialIds, err := mosaicPartialService.FindRepeats(mosaic, maxRepeats)
	if err != nil {
		return err
	}

	if len(macroPartialIds) == 0 {
		return nil
	}

	partialComparisonService := env.MustPartialComparisonService()

	for _, id := range macroPartialIds {
		err = partialComparisonService.DeleteBy("macro_partial_id = ?", id)
		if err != nil {
			return err
		}
	}

	return nil
}
