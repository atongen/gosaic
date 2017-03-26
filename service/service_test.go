package service

import (
	"github.com/atongen/gosaic/model"
)

var (
	serviceFactory ServiceFactory

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

func setTestServiceFactory() {
	var err error
	serviceFactory, err = NewServiceFactory("sqlite3://:memory:")
	if err != nil {
		panic(err)
	}
}
