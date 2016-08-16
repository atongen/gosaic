package service

import (
	"gosaic/model"
	"testing"
)

func setupMosaicPartialServiceTest() (MosaicPartialService, error) {
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

	coverService, err := getTestCoverService(dbMap)
	if err != nil {
		return nil, err
	}

	coverPartialService, err := getTestCoverPartialService(dbMap)
	if err != nil {
		return nil, err
	}

	aspectService, err := getTestAspectService(dbMap)
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

	//partialComparisonService, err := getTestPartialComparisonService(dbMap)
	//if err != nil {
	//	return nil, err
	//}

	mosaicService, err := getTestMosaicService(dbMap)
	if err != nil {
		return nil, err
	}

	mosaicPartialService, err := getTestMosaicPartialService(dbMap)
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

	coverPartials := make([]model.CoverPartial, 6)
	for i := 0; i < 6; i++ {
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
		if i == 6 {
			coverPartial = cp
		} else {
			coverPartials[i] = cp
		}
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

	for i := 0; i < 5; i++ {
		mp, err := macroPartialService.FindOrCreate(&macro, &coverPartials[i])
		if err != nil {
			return nil, err
		}
		if i == 0 {
			macroPartial = *mp
		}
	}

	mosaic = model.Mosaic{
		Name:    "test1",
		MacroId: macro.Id,
	}
	err = mosaicService.Insert(&mosaic)
	if err != nil {
		return nil, err
	}

	return mosaicPartialService, nil
}

func TestMosaicPartialServiceInsert(t *testing.T) {
	mosaicPartialService, err := setupMosaicPartialServiceTest()
	if err != nil {
		t.Fatalf("Unable to setup database: %s\n", err.Error())
	}
	defer mosaicPartialService.DbMap().Db.Close()

	c1 := model.MosaicPartial{
		MosaicId:       mosaic.Id,
		MacroPartialId: macroPartial.Id,
		GidxPartialId:  gidxPartial.Id,
	}

	err = mosaicPartialService.Insert(&c1)
	if err != nil {
		t.Fatalf("Error inserting mosaic partial: %s\n", err.Error())
	}

	if c1.Id == int64(0) {
		t.Fatalf("Inserted mosaic paritial id not set")
	}

	c2, err := mosaicPartialService.Get(c1.Id)
	if err != nil {
		t.Fatalf("Error getting inserted mosaic partial: %s\n", err.Error())
	} else if c2 == nil {
		t.Fatalf("Mosaic partial not inserted\n")
	}

	if c1.Id != c2.Id ||
		c1.MacroPartialId != c2.MacroPartialId ||
		c1.GidxPartialId != c2.GidxPartialId {
		t.Fatalf("Inserted mosaic partial (%+v) does not match: %+v\n", c2, c1)
	}
}
