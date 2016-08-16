package controller

import (
	"gosaic/environment"
	"gosaic/model"
	"gosaic/service"
	"log"
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

	if maxRepeats > 0 {
		if numGidxs*int64(maxRepeats) < numMacroPartials {
			env.Fatalf("Not enough index images (%d) to fill mosaic (%d) with max repeats set to %d", numGidxs, numMacroPartials, maxRepeats)
		}
	}

	mosaic, err := mosaicService.GetOneBy("macro_id = ? AND name = ?", macroId, name)
	if err != nil {
		env.Fatalf("Error checking for existing mosaic: %s\n", err.Error())
	}

	if mosaic == nil {
		mosaic := model.Mosaic{
			Name:    name,
			MacroId: macro.Id,
		}
		err = mosaicService.Insert(&mosaic)
		if err != nil {
			env.Fatalf("Error creating mosaic: %s\n", err.Error())
		}
	}

	createMosaicPartials(env.Log(), mosaic)
}

// TODO: service tests
func createMosaicPartials(l *log.Logger, mosaicPartialService service.MosaicPartialService, mosaic *model.Mosaic) {
	for {
		macroPartial, err := mosaicPartialService.GetRandomMissing(mosaic)
		if macroPartial == nil {
			break
		}

		gidxPartialId, err := partialComparisonService.GetClosest(mosaic, macroPartial, maxRepeats)
		if err != nil {
			l.Fatalf("Error finding closest index image: %s\n", err.Error())
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
