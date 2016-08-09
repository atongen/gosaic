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

	aspect = model.Aspect{Columns: 239, Rows: 170}
	err = aspectService.Insert(&aspect)
	if err != nil {
		return nil, err
	}

	gidx = model.Gidx{
		AspectId:    aspect.Id,
		Path:        "testdata/shaq_bill.jpg",
		Md5sum:      "394c43174e42e043e7b9049e1bb10a39",
		Width:       uint(478),
		Height:      uint(340),
		Orientation: 1,
	}
	err = gidxService.Insert(&gidx)
	if err != nil {
		return nil, err
	}

	gidx2 := model.Gidx{
		AspectId:    aspect.Id,
		Path:        "testdata/eagle.jpg",
		Md5sum:      "5a19b84638fc471d8ec4167ea4e659fb",
		Width:       uint(512),
		Height:      uint(364),
		Orientation: 1,
	}
	err = gidxService.Insert(&gidx2)
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

	_, err = gidxPartialService.FindOrCreate(&gidx2, &aspect)
	if err != nil {
		return nil, err
	}

	coverPartials := make([]model.CoverPartial, 5)
	for i := 0; i < 5; i++ {
		cp := model.CoverPartial{
			CoverId:  cover.Id,
			AspectId: aspect.Id,
			X1:       int64(i),
			Y1:       int64(i),
			X2:       int64(i + 1),
			Y2:       int64(i + 1),
		}
		err = coverPartialService.Insert(&cp)
		if err != nil {
			return nil, err
		}
		coverPartials[i] = cp
	}
	coverPartial = coverPartials[0]

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

	for i := 0; i < 5; i++ {
		mp, err := macroPartialService.FindOrCreate(&macro, &coverPartials[i])
		if err != nil {
			return nil, err
		}
		if i == 0 {
			macroPartial = *mp
		}
	}

	return partialComparisonService, nil
}

func TestPartialComparisonServiceCountMissing(t *testing.T) {
	partialComparisonService, err := setupPartialComparisonServiceTest()
	if err != nil {
		t.Fatalf("Unable to setup database: %s\n", err.Error())
	}
	defer partialComparisonService.DbMap().Db.Close()

	num, err := partialComparisonService.CountMissing(&macro)
	if err != nil {
		t.Fatalf("Error counting missing partial comparisons: %s\n", err.Error())
	}

	if num != 10 {
		t.Fatalf("Expected 10 missing partial comparisons, got %d\n", num)
	}
}

func TestPartialComparisonServiceFindMissing(t *testing.T) {
	partialComparisonService, err := setupPartialComparisonServiceTest()
	if err != nil {
		t.Fatalf("Unable to setup database: %s\n", err.Error())
	}
	defer partialComparisonService.DbMap().Db.Close()

	macroGidxViews, err := partialComparisonService.FindMissing(&macro, 1000)
	if err != nil {
		t.Fatalf("Error finding missing partial comparisons: %s\n", err.Error())
	}

	if len(macroGidxViews) != 10 {
		t.Fatalf("Expected 10 missing partial comparisons, got %d\n", len(macroGidxViews))
	}
}
