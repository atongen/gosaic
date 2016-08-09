package service

import (
	"testing"

	"gosaic/model"

	_ "github.com/mattn/go-sqlite3"
)

func setupPartialComparisonServiceTest() (PartialComparisonService, error) {
	dbMap, err := getTestDbMap()
	if err != nil {
		return nil, err
	}

	gidxService, err := getTestGidxService(dbMap)
	if err != nil {
		return nil, err
	}

	gidxPartialService, err := getTestGidxPartialService(dbMap)
	if err != nil {
		return nil, err
	}

	aspectService, err := getTestAspectService(dbMap)
	if err != nil {
		return nil, err
	}

	coverService, err := getTestCoverService(dbMap)
	if err != nil {
		return nil, err
	}

	coverPartialService, err := getTestCoverPartialService(dbMap)
	if err != nil {
		return nil, err
	}

	macroService, err := getTestMacroService(dbMap)
	if err != nil {
		return nil, err
	}

	macroPartialService, err := getTestMacroPartialService(dbMap)
	if err != nil {
		return nil, err
	}

	partialComparisonService, err := getTestPartialComparisonService(dbMap)
	if err != nil {
		return nil, err
	}

	aspect = model.Aspect{Columns: 87, Rows: 128}
	err = aspectService.Insert(&aspect)
	if err != nil {
		return nil, err
	}

	gidx = model.Gidx{
		AspectId:    aspect.Id,
		Path:        "testdata/matterhorn.jpg",
		Md5sum:      "fcaadee574094a3ae04c6badbbb9ee5e",
		Width:       uint(696),
		Height:      uint(1024),
		Orientation: 1,
	}
	err = gidxService.Insert(&gidx)
	if err != nil {
		return nil, err
	}

	cover = model.Cover{Name: "test1", AspectId: aspect.Id, Width: 1, Height: 1}
	err = coverService.Insert(&cover)
	if err != nil {
		return nil, err
	}

	gp, err := gidxPartialService.FindOrCreate(&gidx, &aspect)
	if err != nil {
		return nil, err
	}
	gidxPartial = *gp

	coverPartial = model.CoverPartial{
		CoverId:  cover.Id,
		AspectId: aspect.Id,
		X1:       0,
		Y1:       0,
		X2:       1,
		Y2:       1,
	}
	err = coverPartialService.Insert(&coverPartial)
	if err != nil {
		return nil, err
	}

	macro = model.Macro{
		CoverId:     cover.Id,
		AspectId:    aspect.Id,
		Path:        "testdata/matterhorn.jpg",
		Md5sum:      "fcaadee574094a3ae04c6badbbb9ee5e",
		Width:       uint(696),
		Height:      uint(1024),
		Orientation: 1,
	}
	err = macroService.Insert(&macro)
	if err != nil {
		return nil, err
	}

	mp, err := macroPartialService.FindOrCreate(&macro, &coverPartial)
	if err != nil {
		return nil, err
	}
	macroPartial = *mp

	return partialComparisonService, nil
}

func TestPartialComparisonServiceFindMissing(t *testing.T) {
	partialComparisonService, err := setupPartialComparisonServiceTest()
	if err != nil {
		t.Fatalf("Unable to setup database: %s\n", err.Error())
	}
	defer partialComparisonService.DbMap().Db.Close()

	partialComparisons, err := partialComparisonService.FindMissing(&macro, 1000)
	if err != nil {
		t.Fatalf("Error finding missing partial comparisons: %s\n", err.Error())
	}

	t.Fatalf("missing partial comparisons: %d\n", len(partialComparisons))
}
