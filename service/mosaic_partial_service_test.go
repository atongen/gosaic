package service

import (
	"testing"

	"github.com/atongen/gosaic/model"
)

func setupMosaicPartialServiceTest() {
	setTestServiceFactory()
	gidxService := serviceFactory.MustGidxService()
	gidxPartialService := serviceFactory.MustGidxPartialService()
	coverService := serviceFactory.MustCoverService()
	coverPartialService := serviceFactory.MustCoverPartialService()
	aspectService := serviceFactory.MustAspectService()
	macroService := serviceFactory.MustMacroService()
	macroPartialService := serviceFactory.MustMacroPartialService()
	mosaicService := serviceFactory.MustMosaicService()

	aspect = model.Aspect{Columns: 87, Rows: 128}
	err := aspectService.Insert(&aspect)
	if err != nil {
		panic(err)
	}

	gidx = model.Gidx{
		AspectId:    aspect.Id,
		Path:        "testdata/shaq_bill.jpg",
		Md5sum:      "394c43174e42e043e7b9049e1bb10a39",
		Width:       478,
		Height:      340,
		Orientation: 1,
	}
	err = gidxService.Insert(&gidx)
	if err != nil {
		panic(err)
	}

	gidx2 := model.Gidx{
		AspectId:    aspect.Id,
		Path:        "testdata/eagle.jpg",
		Md5sum:      "5a19b84638fc471d8ec4167ea4e659fb",
		Width:       512,
		Height:      364,
		Orientation: 1,
	}
	err = gidxService.Insert(&gidx2)
	if err != nil {
		panic(err)
	}

	cover = model.Cover{AspectId: aspect.Id, Width: 1, Height: 1}
	err = coverService.Insert(&cover)
	if err != nil {
		panic(err)
	}

	gp, err := gidxPartialService.FindOrCreate(&gidx, &aspect)
	if err != nil {
		panic(err)
	}
	gidxPartial = *gp

	_, err = gidxPartialService.FindOrCreate(&gidx2, &aspect)
	if err != nil {
		panic(err)
	}

	coverPartials := make([]model.CoverPartial, 6)
	for i := 0; i < 6; i++ {
		cp := model.CoverPartial{
			CoverId:  cover.Id,
			AspectId: aspect.Id,
			X1:       i,
			Y1:       i,
			X2:       i + 1,
			Y2:       i + 1,
		}
		err = coverPartialService.Insert(&cp)
		if err != nil {
			panic(err)
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
		Width:       696,
		Height:      1024,
		Orientation: 1,
	}
	err = macroService.Insert(&macro)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 5; i++ {
		mp, err := macroPartialService.FindOrCreate(&macro, &coverPartials[i])
		if err != nil {
			panic(err)
		}
		if i == 0 {
			macroPartial = *mp
		}
	}

	mosaic = model.Mosaic{
		MacroId: macro.Id,
	}
	err = mosaicService.Insert(&mosaic)
	if err != nil {
		panic(err)
	}
}

func TestMosaicPartialServiceInsert(t *testing.T) {
	setupMosaicPartialServiceTest()
	mosaicPartialService := serviceFactory.MustMosaicPartialService()
	defer mosaicPartialService.Close()

	c1 := model.MosaicPartial{
		MosaicId:       mosaic.Id,
		MacroPartialId: macroPartial.Id,
		GidxPartialId:  gidxPartial.Id,
	}

	err := mosaicPartialService.Insert(&c1)
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

func TestMosaicPartialServiceCount(t *testing.T) {
	setupMosaicPartialServiceTest()
	mosaicPartialService := serviceFactory.MustMosaicPartialService()
	defer mosaicPartialService.Close()

	num, err := mosaicPartialService.Count(&mosaic)
	if err != nil {
		t.Fatalf("Error counting missing mosaic partials: %s\n", err.Error())
	}

	if num != 0 {
		t.Fatalf("Expected 5 mosaic partials, got %d\n", num)
	}

	c1 := model.MosaicPartial{
		MosaicId:       mosaic.Id,
		MacroPartialId: macroPartial.Id,
		GidxPartialId:  gidxPartial.Id,
	}

	err = mosaicPartialService.Insert(&c1)
	if err != nil {
		t.Fatalf("Error inserting mosaic partial: %s\n", err.Error())
	}

	num, err = mosaicPartialService.Count(&mosaic)
	if err != nil {
		t.Fatalf("Error counting missing mosaic partials: %s\n", err.Error())
	}

	if num != 1 {
		t.Fatalf("Expected 4 missing mosaic partials, got %d\n", num)
	}
}

func TestMosaicPartialServiceCountMissing(t *testing.T) {
	setupMosaicPartialServiceTest()
	mosaicPartialService := serviceFactory.MustMosaicPartialService()
	defer mosaicPartialService.Close()

	num, err := mosaicPartialService.CountMissing(&mosaic)
	if err != nil {
		t.Fatalf("Error counting missing mosaic partials: %s\n", err.Error())
	}

	if num != 5 {
		t.Fatalf("Expected 5 missing mosaic partials, got %d\n", num)
	}

	c1 := model.MosaicPartial{
		MosaicId:       mosaic.Id,
		MacroPartialId: macroPartial.Id,
		GidxPartialId:  gidxPartial.Id,
	}

	err = mosaicPartialService.Insert(&c1)
	if err != nil {
		t.Fatalf("Error inserting mosaic partial: %s\n", err.Error())
	}

	num, err = mosaicPartialService.CountMissing(&mosaic)
	if err != nil {
		t.Fatalf("Error counting missing mosaic partials: %s\n", err.Error())
	}

	if num != 4 {
		t.Fatalf("Expected 4 missing mosaic partials, got %d\n", num)
	}
}

func TestMosaicPartialServiceGetMissing(t *testing.T) {
	setupMosaicPartialServiceTest()
	mosaicPartialService := serviceFactory.MustMosaicPartialService()
	defer mosaicPartialService.Close()

	mp1, err := mosaicPartialService.GetMissing(&mosaic)
	if err != nil {
		t.Fatalf("Error finding missing mosaic partial: %s\n", err.Error())
	}
	if mp1 == nil {
		t.Fatal("Missing macro partial not found")
	}

	if mp1.Id != int64(1) {
		t.Fatalf("Expected first missing macro partial to have id 1, got %d\n", mp1.Id)
	}

	c1 := model.MosaicPartial{
		MosaicId:       mosaic.Id,
		MacroPartialId: mp1.Id,
		GidxPartialId:  gidxPartial.Id,
	}

	err = mosaicPartialService.Insert(&c1)
	if err != nil {
		t.Fatalf("Error inserting mosaic partial: %s\n", err.Error())
	}

	c2, err := mosaicPartialService.GetMissing(&mosaic)
	if err != nil {
		t.Fatalf("Error finding missing mosaic partial: %s\n", err.Error())
	}
	if c2 == nil {
		t.Fatal("Missing macro partial not found")
	}

	if c2.Id != int64(2) {
		t.Fatalf("Expected second missing macro partial to have id 2, got %d\n", c2.Id)
	}
}

func TestMosaicPartialServiceGetRandomMissing(t *testing.T) {
	setupMosaicPartialServiceTest()
	mosaicPartialService := serviceFactory.MustMosaicPartialService()
	defer mosaicPartialService.Close()

	mp1, err := mosaicPartialService.GetRandomMissing(&mosaic)
	if err != nil {
		t.Fatalf("Error finding random missing mosaic partial: %s\n", err.Error())
	}
	if mp1 == nil {
		t.Fatal("Missing macro partial not found")
	}

	c1 := model.MosaicPartial{
		MosaicId:       mosaic.Id,
		MacroPartialId: mp1.Id,
		GidxPartialId:  gidxPartial.Id,
	}

	err = mosaicPartialService.Insert(&c1)
	if err != nil {
		t.Fatalf("Error inserting mosaic partial: %s\n", err.Error())
	}

	for i := 0; i < 10; i++ {
		c2, err := mosaicPartialService.GetRandomMissing(&mosaic)
		if err != nil {
			t.Fatalf("Error finding random missing mosaic partial: %s\n", err.Error())
		}
		if c2 == nil {
			t.Fatal("Missing macro partial not found")
		}

		if c2.Id == mp1.Id {
			t.Fatalf("Random missing macro partial has same id as inserted: %d\n", c2.Id)
		}
	}
}

func TestMosaicPartialServiceFindAllPartialViews(t *testing.T) {
	setupMosaicPartialServiceTest()
	mosaicPartialService := serviceFactory.MustMosaicPartialService()
	defer mosaicPartialService.Close()

	c1 := model.MosaicPartial{
		MosaicId:       mosaic.Id,
		MacroPartialId: macroPartial.Id,
		GidxPartialId:  gidxPartial.Id,
	}

	err := mosaicPartialService.Insert(&c1)
	if err != nil {
		t.Fatalf("Error inserting mosaic partial: %s\n", err.Error())
	}

	if c1.Id == int64(0) {
		t.Fatalf("Inserted mosaic paritial id not set")
	}

	views, err := mosaicPartialService.FindAllPartialViews(&mosaic, "mosaic_partials.id asc", 1000, 0)
	if err != nil {
		t.Fatalf("Error finding all mosaic partial views: %s\n", err.Error())
	}

	if len(views) != 1 {
		t.Fatalf("Expected 1 mosaic partial view, got: %d\n", len(views))
	}

	view := views[0]
	if view.MosaicPartialId != c1.Id {
		t.Fatalf("Expected mosaic partial id %d, got: %d\n", c1.Id, view.MosaicPartialId)
	}

	if gidx.Id != view.Gidx.Id ||
		gidx.AspectId != view.Gidx.AspectId ||
		gidx.Md5sum != view.Gidx.Md5sum ||
		gidx.Path != view.Gidx.Path ||
		gidx.Width != view.Gidx.Width ||
		gidx.Height != view.Gidx.Height ||
		gidx.Orientation != view.Gidx.Orientation {
		t.Error("Found view gidx does not match data")
	}

	if coverPartial.Id != view.CoverPartial.Id ||
		coverPartial.CoverId != view.CoverPartial.CoverId ||
		coverPartial.AspectId != view.CoverPartial.AspectId ||
		coverPartial.X1 != view.CoverPartial.X1 ||
		coverPartial.Y1 != view.CoverPartial.Y1 ||
		coverPartial.X2 != view.CoverPartial.X2 ||
		coverPartial.Y2 != view.CoverPartial.Y2 {
		t.Fatalf("Inserted cover partial (%+v) does not match: %+v\n", view.CoverPartial, coverPartial)
	}
}

func TestMosaicPartialServiceFindRepeats(t *testing.T) {
	setupMosaicPartialServiceTest()
	mosaicPartialService := serviceFactory.MustMosaicPartialService()
	defer mosaicPartialService.Close()

	c1 := model.MosaicPartial{
		MosaicId:       mosaic.Id,
		MacroPartialId: macroPartial.Id,
		GidxPartialId:  gidxPartial.Id,
	}

	err := mosaicPartialService.Insert(&c1)
	if err != nil {
		t.Fatalf("Error inserting mosaic partial: %s\n", err.Error())
	}

	ids, err := mosaicPartialService.FindRepeats(&mosaic, 2)
	if err != nil {
		t.Fatalf("Error finding macro partials with repeats: %s\n", err.Error())
	}

	num := len(ids)
	if num != 0 {
		t.Fatalf("Expected 0 gidx partials with 2 or more macro partial duplicates, but got %d\n", num)
	}

	c2 := model.MosaicPartial{
		MosaicId:       mosaic.Id,
		MacroPartialId: macroPartial.Id + 1,
		GidxPartialId:  gidxPartial.Id,
	}

	err = mosaicPartialService.Insert(&c2)
	if err != nil {
		t.Fatalf("Error inserting mosaic partial: %s\n", err.Error())
	}

	ids, err = mosaicPartialService.FindRepeats(&mosaic, 2)
	if err != nil {
		t.Fatal("Error finding macro partials with repeats")
	}

	num = len(ids)
	if num != 1 {
		t.Fatalf("Expected 1 gidx partials with 2 or more macro partial duplicates, but got %d\n", num)
	}
}
