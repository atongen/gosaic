package service

import (
	"database/sql"
	"gosaic/database"
	"gosaic/model"
	"sync"

	"gopkg.in/gorp.v1"
)

var (
	m           sync.Mutex
	cachedDbMap *gorp.DbMap

	// TODO: replace w/ factory
	gidx          model.Gidx
	gidxPartial   model.GidxPartial
	aspect        model.Aspect
	cover         model.Cover
	coverPartial  model.CoverPartial
	macro         model.Macro
	macroPartial  model.MacroPartial
	mosaic        model.Mosaic
	mosaicPartial model.MosaicPartial
)

func setTestDbMap() {
	m.Lock()
	defer m.Unlock()

	_, err := buildTestDb()
	if err != nil {
		panic(err)
	}
}

func getTestDbMap() {
	m.Lock()
	defer m.Unlock()

	if cachedDbMap == nil {
		_, err := buildTestDbMap()
		if err != nil {
			panic(err)
		}
	}
}

func buildTestDbMap() error {
	cachedDbMap = nil

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return err
	}

	_, err = database.Migrate(db)
	if err != nil {
		return err
	}

	cachedDbMap = &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}

	return nil
}

func getTestGidxService() GidxService {
	gidxService := NewGidxService(getTestDbMap())
	err := gidxService.Register()
	if err != nil {
		return panic(err)
	}

	return gidxService
}

func getTestAspectService() AspectService {
	aspectService := NewAspectService(getTestDbMap())
	err := aspectService.Register()
	if err != nil {
		panic(err)
	}

	return aspectService
}

func getTestGidxPartialService() GidxPartialService {
	gidxPartialService := NewGidxPartialService(getTestDbMap())
	err := gidxPartialService.Register()
	if err != nil {
		panic(err)
	}

	return gidxPartialService
}

func getTestCoverService() CoverService {
	coverService := NewCoverService(getTestDbMap())
	err := coverService.Register()
	if err != nil {
		panic(err)
	}

	return coverService
}

func getTestCoverPartialService() CoverPartialService {
	coverPartialService := NewCoverPartialService(getTestDbMap())
	err := coverPartialService.Register()
	if err != nil {
		panic(err)
	}

	return coverPartialService
}

func getTestMacroService() MacroService {
	macroService := NewMacroService(getTestDbMap())
	err := macroService.Register()
	if err != nil {
		panic(err)
	}

	return macroService
}

func getTestMacroPartialService() MacroPartialService {
	macroPartialService := NewMacroPartialService(getTestDbMap())
	err := macroPartialService.Register()
	if err != nil {
		panic(err)
	}

	return macroPartialService
}

func getTestPartialComparisonService() PartialComparisonService {
	partialComparisonService := NewPartialComparisonService(getTestDbMap())
	err := partialComparisonService.Register()
	if err != nil {
		panic(err)
	}

	return partialComparisonService
}

func getTestMosaicService() MosaicService {
	mosaicService := NewMosaicService(getTestDbMap())
	err := mosaicService.Register()
	if err != nil {
		panic(err)
	}

	return mosaicService
}

func getTestMosaicPartialService() MosaicPartialService {
	mosaicPartialService := NewMosaicPartialService(getTestDbMap())
	err := mosaicPartialService.Register()
	if err != nil {
		panic(err)
	}

	return mosaicPartialService
}
