package controller

import (
	"gosaic/environment"
	"gosaic/model"
	"gosaic/service"
	"log"
	"math"
)

func MosaicBuild(env environment.Environment, name string, macroId int64, maxRepeats int) {
	gidxService, err := env.GidxService()
	if err != nil {
		env.Fatalf("Error getting index service: %s\n", err.Error())
	}

	macroService, err := env.MacroService()
	if err != nil {
		env.Fatalf("Error getting macro service: %s\n", err.Error())
	}

	macroPartialService, err := env.MacroPartialService()
	if err != nil {
		env.Fatalf("Error getting macro partial service: %s\n", err.Error())
	}

	mosaicService, err := env.MosaicService()
	if err != nil {
		env.Fatalf("Error getting mosaic service: %s\n", err.Error())
	}

	mosaicPartialService, err := env.MosaicPartialService()
	if err != nil {
		env.Fatalf("Error getting mosaic partial service: %s\n", err.Error())
	}

	partialComparisonService, err := env.PartialComparisonService()
	if err != nil {
		env.Fatalf("Error getting partial comparison service: %s\n", err.Error())
	}

	macro, err := macroService.Get(macroId)
	if err != nil {
		env.Fatalf("Error getting macro: %s\n", err.Error())
	}

	if macro == nil {
		env.Fatalf("Macro %d not found\n", macroId)
	}

	numMacroPartials, err := macroPartialService.CountBy("macro_id", macro.Id)
	if err != nil {
		env.Fatalf("Error counting macro partials: %s\n", err.Error())
	}

	numGidxs, err := gidxService.Count()
	if err != nil {
		env.Fatalf("Error counting index images: %s\n", err.Error())
	}

	// maxRepeats == 0 is unrestricted
	// maxRepeats > 0 sets explicitly
	// maxRepeats == -1 (<0) calculates minimum
	if maxRepeats < 0 {
		if numGidxs > numMacroPartials {
			maxRepeats = 1
		} else {
			maxRepeats = int(math.Ceil(float64(numMacroPartials) / float64(numGidxs)))
		}
	} else if maxRepeats > 0 {
		if numGidxs*int64(maxRepeats) < numMacroPartials {
			env.Fatalf("Not enough index images (%d) to fill mosaic (%d) with max repeats set to %d", numGidxs, numMacroPartials, maxRepeats)
		}
	}

	mosaic, err := mosaicService.GetOneBy("macro_id = ? AND name = ?", macroId, name)
	if err != nil {
		env.Fatalf("Error checking for existing mosaic: %s\n", err.Error())
	}

	if mosaic == nil {
		mosaic = &model.Mosaic{
			Name:    name,
			MacroId: macro.Id,
		}
		err = mosaicService.Insert(mosaic)
		if err != nil {
			env.Fatalf("Error creating mosaic: %s\n", err.Error())
		}
	}

	env.Printf("Creating mosaic with %d total partials\n", numMacroPartials)
	createMosaicPartials(env.Log(), mosaicPartialService, partialComparisonService, mosaic, maxRepeats)
}

func createMosaicPartials(l *log.Logger, mosaicPartialService service.MosaicPartialService, partialComparisonService service.PartialComparisonService, mosaic *model.Mosaic, maxRepeats int) {
	numMissing, err := mosaicPartialService.CountMissing(mosaic)
	if err != nil {
		l.Fatalf("Error counting missing mosaic partials: %s\n", err.Error())
	}

	if numMissing == 0 {
		// we are done
		return
	}

	l.Printf("Building %d missing mosaic partials\n", numMissing)

	for {
		macroPartial := mosaicPartialService.GetRandomMissing(mosaic)
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
			l.Fatalf("Error finding closest index image: %s\n", err.Error())
		}

		if gidxPartialId == int64(0) {
			l.Fatal("Unable to find index to fill mosaic")
		}

		mosaicPartial := model.MosaicPartial{
			MosaicId:       mosaic.Id,
			MacroPartialId: macroPartial.Id,
			GidxPartialId:  gidxPartialId,
		}
		err = mosaicPartialService.Insert(&mosaicPartial)
		if err != nil {
			l.Fatalf("Error inserting mosaic partial: %s\n", err.Error())
		}
	}
}
