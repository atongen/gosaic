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

	cachedAspectService            AspectService
	cachedGidxService              GidxService
	cachedGidxPartialService       GidxPartialService
	cachedCoverService             CoverService
	cachedCoverPartialService      CoverPartialService
	cachedMacroService             MacroService
	cachedMacroPartialService      MacroPartialService
	cachedMosaicService            MosaicService
	cachedMosaicPartialService     MosaicPartialService
	cachedPartialComparisonService PartialComparisonService
	cachedQuadDistService          QuadDistService
	cachedProjectService           ProjectService

	aspect        model.Aspect
	gidx          model.Gidx
	gidxPartial   model.GidxPartial
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

	err := _buildTestDbMap()
	if err != nil {
		panic(err)
	}
}

func getTestAspectService() AspectService {
	m.Lock()
	defer m.Unlock()

	if cachedAspectService == nil {
		cachedAspectService = NewAspectService(_getTestDbMap())
		err := cachedAspectService.Register()
		if err != nil {
			panic(err)
		}
	}

	return cachedAspectService
}

func getTestGidxService() GidxService {
	m.Lock()
	defer m.Unlock()

	if cachedGidxService == nil {
		cachedGidxService = NewGidxService(_getTestDbMap())
		err := cachedGidxService.Register()
		if err != nil {
			panic(err)
		}
	}

	return cachedGidxService
}

func getTestGidxPartialService() GidxPartialService {
	m.Lock()
	defer m.Unlock()

	if cachedGidxPartialService == nil {
		cachedGidxPartialService = NewGidxPartialService(_getTestDbMap())
		err := cachedGidxPartialService.Register()
		if err != nil {
			panic(err)
		}
	}

	return cachedGidxPartialService
}

func getTestCoverService() CoverService {
	m.Lock()
	defer m.Unlock()

	if cachedCoverService == nil {
		cachedCoverService = NewCoverService(_getTestDbMap())
		err := cachedCoverService.Register()
		if err != nil {
			panic(err)
		}
	}

	return cachedCoverService
}

func getTestCoverPartialService() CoverPartialService {
	m.Lock()
	defer m.Unlock()

	if cachedCoverPartialService == nil {
		cachedCoverPartialService = NewCoverPartialService(_getTestDbMap())
		err := cachedCoverPartialService.Register()
		if err != nil {
			panic(err)
		}
	}

	return cachedCoverPartialService
}

func getTestMacroService() MacroService {
	m.Lock()
	defer m.Unlock()

	if cachedMacroService == nil {
		cachedMacroService = NewMacroService(_getTestDbMap())
		err := cachedMacroService.Register()
		if err != nil {
			panic(err)
		}
	}

	return cachedMacroService
}

func getTestMacroPartialService() MacroPartialService {
	m.Lock()
	defer m.Unlock()

	if cachedMacroPartialService == nil {
		cachedMacroPartialService = NewMacroPartialService(_getTestDbMap())
		err := cachedMacroPartialService.Register()
		if err != nil {
			panic(err)
		}
	}

	return cachedMacroPartialService
}

func getTestPartialComparisonService() PartialComparisonService {
	m.Lock()
	defer m.Unlock()

	if cachedPartialComparisonService == nil {
		cachedPartialComparisonService = NewPartialComparisonService(_getTestDbMap())
		err := cachedPartialComparisonService.Register()
		if err != nil {
			panic(err)
		}
	}

	return cachedPartialComparisonService
}

func getTestMosaicService() MosaicService {
	m.Lock()
	defer m.Unlock()

	if cachedMosaicService == nil {
		cachedMosaicService = NewMosaicService(_getTestDbMap())
		err := cachedMosaicService.Register()
		if err != nil {
			panic(err)
		}
	}

	return cachedMosaicService
}

func getTestMosaicPartialService() MosaicPartialService {
	m.Lock()
	defer m.Unlock()

	if cachedMosaicPartialService == nil {
		cachedMosaicPartialService = NewMosaicPartialService(_getTestDbMap())
		err := cachedMosaicPartialService.Register()
		if err != nil {
			panic(err)
		}
	}

	return cachedMosaicPartialService
}

func getTestQuadDistService() QuadDistService {
	m.Lock()
	defer m.Unlock()

	if cachedQuadDistService == nil {
		cachedQuadDistService = NewQuadDistService(_getTestDbMap())
		err := cachedQuadDistService.Register()
		if err != nil {
			panic(err)
		}
	}

	return cachedQuadDistService
}

func getTestProjectService() ProjectService {
	m.Lock()
	defer m.Unlock()

	if cachedProjectService == nil {
		cachedProjectService = NewProjectService(_getTestDbMap())
		err := cachedProjectService.Register()
		if err != nil {
			panic(err)
		}
	}

	return cachedProjectService
}

func _getTestDbMap() *gorp.DbMap {
	if cachedDbMap == nil {
		err := _buildTestDbMap()
		if err != nil {
			panic(err)
		}
	}
	return cachedDbMap
}

func _buildTestDbMap() error {
	_resetTestDbMap()

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

func _resetTestDbMap() {
	cachedDbMap = nil
	cachedAspectService = nil
	cachedGidxService = nil
	cachedGidxPartialService = nil
	cachedCoverService = nil
	cachedCoverPartialService = nil
	cachedMacroService = nil
	cachedMacroPartialService = nil
	cachedMosaicService = nil
	cachedMosaicPartialService = nil
	cachedPartialComparisonService = nil
	cachedQuadDistService = nil
	cachedProjectService = nil
}
