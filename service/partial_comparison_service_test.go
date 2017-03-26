package service

import (
	"testing"

	"github.com/atongen/gosaic/model"

	_ "github.com/mattn/go-sqlite3"
)

func setupPartialComparisonServiceTest() {
	setTestServiceFactory()
	gidxService := serviceFactory.MustGidxService()
	gidxPartialService := serviceFactory.MustGidxPartialService()
	aspectService := serviceFactory.MustAspectService()
	coverService := serviceFactory.MustCoverService()
	coverPartialService := serviceFactory.MustCoverPartialService()
	macroService := serviceFactory.MustMacroService()
	macroPartialService := serviceFactory.MustMacroPartialService()

	aspect = model.Aspect{Columns: 239, Rows: 170}
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
}

func TestPartialComparisonServiceInsert(t *testing.T) {
	setupPartialComparisonServiceTest()
	partialComparisonService := serviceFactory.MustPartialComparisonService()
	defer partialComparisonService.Close()

	pc := model.PartialComparison{
		MacroPartialId: macroPartial.Id,
		GidxPartialId:  gidxPartial.Id,
		Dist:           0.5,
	}

	err := partialComparisonService.Insert(&pc)
	if err != nil {
		t.Fatalf("Error inserting partial comparison: %s\n", err.Error())
	}

	if pc.Id == int64(0) {
		t.Fatalf("Inserted partial comparison id not set")
	}

	pc2, err := partialComparisonService.Get(pc.Id)
	if err != nil {
		t.Fatalf("Error getting inserted partial comparison: %s\n", err.Error())
	} else if pc2 == nil {
		t.Fatalf("Partial comparison not inserted\n")
	}

	if pc.Id != pc2.Id ||
		pc.MacroPartialId != pc2.MacroPartialId ||
		pc.GidxPartialId != pc2.GidxPartialId ||
		pc.Dist != pc2.Dist {
		t.Fatal("Inserted macro partial data does not match")
	}
}

func TestPartialComparisonServiceBulkInsert(t *testing.T) {
	setupPartialComparisonServiceTest()
	partialComparisonService := serviceFactory.MustPartialComparisonService()
	defer partialComparisonService.Close()

	macroGidxViews, err := partialComparisonService.FindMissing(&macro, 1000)
	if err != nil {
		t.Fatalf("Error finding missing partial comparisons: %s\n", err.Error())
	}

	if len(macroGidxViews) < 2 {
		t.Fatalf("Expected more than 1 missing partial comparisons, got %d", len(macroGidxViews))
	}

	partialComparisons := make([]*model.PartialComparison, len(macroGidxViews))
	for i, mgv := range macroGidxViews {
		partialComparisons[i] = &model.PartialComparison{
			MacroPartialId: mgv.MacroPartial.Id,
			GidxPartialId:  mgv.GidxPartial.Id,
			Dist:           0.5,
		}
	}

	num, err := partialComparisonService.BulkInsert(partialComparisons)
	if err != nil {
		t.Fatalf("Error bulk inserting partial comparisons: %s\n", err.Error())
	}

	if num != int64(len(macroGidxViews)) {
		t.Fatalf("Expected %d affected rows for bulk insert, got %d\n", len(macroGidxViews), num)
	}

	count, err := partialComparisonService.CountMissing(&macro)
	if err != nil {
		t.Fatalf("Error counting missing partial comparisons: %s\n", err.Error())
	}

	if count != 0 {
		t.Fatalf("Expected 0 missing partial comparisons, got %d\n", count)
	}

}

func TestPartialComparisonServiceUpdate(t *testing.T) {
	setupPartialComparisonServiceTest()
	partialComparisonService := serviceFactory.MustPartialComparisonService()
	defer partialComparisonService.Close()

	pc := model.PartialComparison{
		MacroPartialId: macroPartial.Id,
		GidxPartialId:  gidxPartial.Id,
		Dist:           0.5,
	}

	err := partialComparisonService.Insert(&pc)
	if err != nil {
		t.Fatalf("Error inserting partial comparison: %s\n", err.Error())
	}

	pc.Dist = 0.24
	err = partialComparisonService.Update(&pc)
	if err != nil {
		t.Fatalf("Error updating partial comparison: %s\n", err.Error())
	}

	pc2, err := partialComparisonService.Get(pc.Id)
	if err != nil {
		t.Fatalf("Error getting updated partial comparison: %s\n", err.Error())
	} else if pc2 == nil {
		t.Fatalf("Partial comparison not inserted\n")
	}

	if pc2.Dist != 0.24 {
		t.Fatal("Updated partial comparison data does not match")
	}
}

func TestPartialComparisonServiceDelete(t *testing.T) {
	setupPartialComparisonServiceTest()
	partialComparisonService := serviceFactory.MustPartialComparisonService()
	defer partialComparisonService.Close()

	pc := model.PartialComparison{
		MacroPartialId: macroPartial.Id,
		GidxPartialId:  gidxPartial.Id,
		Dist:           0.5,
	}

	err := partialComparisonService.Insert(&pc)
	if err != nil {
		t.Fatalf("Error inserting partial comparison: %s\n", err.Error())
	}

	err = partialComparisonService.Delete(&pc)
	if err != nil {
		t.Fatalf("Error deleting partial comparison: %s\n", err.Error())
	}

	pc2, err := partialComparisonService.Get(pc.Id)
	if err != nil {
		t.Fatalf("Error getting deleted partial comparison: %s\n", err.Error())
	} else if pc2 != nil {
		t.Fatalf("partial comparison not deleted\n")
	}
}

func TestPartialComparisonServiceDeleteBy(t *testing.T) {
	setupPartialComparisonServiceTest()
	partialComparisonService := serviceFactory.MustPartialComparisonService()
	defer partialComparisonService.Close()

	pc := model.PartialComparison{
		MacroPartialId: macroPartial.Id,
		GidxPartialId:  gidxPartial.Id,
		Dist:           0.5,
	}

	err := partialComparisonService.Insert(&pc)
	if err != nil {
		t.Fatalf("Error inserting partial comparison: %s\n", err.Error())
	}

	err = partialComparisonService.DeleteBy("macro_partial_id = ?", macroPartial.Id)
	if err != nil {
		t.Fatalf("Error deleting by partial comparison: %s\n", err.Error())
	}

	pc2, err := partialComparisonService.Get(pc.Id)
	if err != nil {
		t.Fatalf("Error getting deleted partial comparison: %s\n", err.Error())
	} else if pc2 != nil {
		t.Fatalf("partial comparison not deleted\n")
	}
}

func TestPartialComparisonServiceDeleteFrom(t *testing.T) {
	setupPartialComparisonServiceTest()
	partialComparisonService := serviceFactory.MustPartialComparisonService()
	defer partialComparisonService.Close()

	pc := model.PartialComparison{
		MacroPartialId: macroPartial.Id,
		GidxPartialId:  gidxPartial.Id,
		Dist:           0.5,
	}

	err := partialComparisonService.Insert(&pc)
	if err != nil {
		t.Fatalf("Error inserting partial comparison: %s\n", err.Error())
	}

	num, err := partialComparisonService.Count()
	if err != nil {
		t.Fatalf("Error counting partial comparison: %s\n", err.Error())
	}

	if num == 0 {
		t.Fatalf("Expected %d partial comparisons got %d\n", 0, num)
	}

	err = partialComparisonService.DeleteFrom(&macro)
	if err != nil {
		t.Fatalf("Error deleting from partial comparison macro: %s\n", err.Error())
	}

	num2, err := partialComparisonService.Count()
	if err != nil {
		t.Fatalf("Error counting partial comparison: %s\n", err.Error())
	}

	if num2 != 0 {
		t.Fatalf("Partial comparisons not deleted. Expected %d, got %d\n", 0, num2)
	}
}

func TestPartialComparisonServiceGetOneBy(t *testing.T) {
	setupPartialComparisonServiceTest()
	partialComparisonService := serviceFactory.MustPartialComparisonService()
	defer partialComparisonService.Close()

	pc := model.PartialComparison{
		MacroPartialId: macroPartial.Id,
		GidxPartialId:  gidxPartial.Id,
		Dist:           0.5,
	}

	err := partialComparisonService.Insert(&pc)
	if err != nil {
		t.Fatalf("Error inserting partial comparison: %s\n", err.Error())
	}

	pc2, err := partialComparisonService.GetOneBy("macro_partial_id = ? and gidx_partial_id = ?", pc.MacroPartialId, pc.GidxPartialId)
	if err != nil {
		t.Fatalf("Error getting one by partial comparison: %s\n", err.Error())
	} else if pc2 == nil {
		t.Fatalf("partial comparison not found by\n")
	}

	if pc2.MacroPartialId != pc.MacroPartialId ||
		pc2.GidxPartialId != pc.GidxPartialId {
		t.Fatal("partial comparison macro id does not match")
	}
}

func TestPartialComparisonServiceGetOneByNot(t *testing.T) {
	setupPartialComparisonServiceTest()
	partialComparisonService := serviceFactory.MustPartialComparisonService()
	defer partialComparisonService.Close()

	_, err := partialComparisonService.GetOneBy("macro_partial_id = ? and gidx_partial_id = ?", macroPartial.Id, gidxPartial.Id)
	if err == nil {
		t.Fatalf("Getting one by partial comparison did not fail")
	}
}

func TestPartialComparisonServiceExistsBy(t *testing.T) {
	setupPartialComparisonServiceTest()
	partialComparisonService := serviceFactory.MustPartialComparisonService()
	defer partialComparisonService.Close()

	pc := model.PartialComparison{
		MacroPartialId: macroPartial.Id,
		GidxPartialId:  gidxPartial.Id,
		Dist:           0.5,
	}

	err := partialComparisonService.Insert(&pc)
	if err != nil {
		t.Fatalf("Error inserting partial comparison: %s\n", err.Error())
	}

	found, err := partialComparisonService.ExistsBy("macro_partial_id = ? and gidx_partial_id = ?", pc.MacroPartialId, pc.GidxPartialId)
	if err != nil {
		t.Fatalf("Error getting one by partial comparison: %s\n", err.Error())
	}

	if !found {
		t.Fatalf("Partial comparison not exists by\n")
	}
}

func TestPartialComparisonServiceExistsByNot(t *testing.T) {
	setupPartialComparisonServiceTest()
	partialComparisonService := serviceFactory.MustPartialComparisonService()
	defer partialComparisonService.Close()

	found, err := partialComparisonService.ExistsBy("macro_partial_id = ? and gidx_partial_id = ?", macroPartial.Id, gidxPartial.Id)
	if err != nil {
		t.Fatalf("Error getting exists by partial comparison: %s\n", err.Error())
	}

	if found {
		t.Fatalf("Partial comparison exists by\n")
	}
}

func TestPartialComparisonServiceCount(t *testing.T) {
	setupPartialComparisonServiceTest()
	partialComparisonService := serviceFactory.MustPartialComparisonService()
	defer partialComparisonService.Close()

	pc := model.PartialComparison{
		MacroPartialId: macroPartial.Id,
		GidxPartialId:  gidxPartial.Id,
		Dist:           0.5,
	}

	err := partialComparisonService.Insert(&pc)
	if err != nil {
		t.Fatalf("Error inserting partial comparison: %s\n", err.Error())
	}

	num, err := partialComparisonService.Count()
	if err != nil {
		t.Fatalf("Error counting partial comparison: %s\n", err.Error())
	}

	if num != int64(1) {
		t.Fatalf("Partial comparison count incorrect\n")
	}
}

func TestPartialComparisonServiceCountBy(t *testing.T) {
	setupPartialComparisonServiceTest()
	partialComparisonService := serviceFactory.MustPartialComparisonService()
	defer partialComparisonService.Close()

	pc := model.PartialComparison{
		MacroPartialId: macroPartial.Id,
		GidxPartialId:  gidxPartial.Id,
		Dist:           0.5,
	}

	err := partialComparisonService.Insert(&pc)
	if err != nil {
		t.Fatalf("Error inserting partial comparison: %s\n", err.Error())
	}

	num, err := partialComparisonService.CountBy("macro_partial_id = ? and gidx_partial_id = ?", pc.MacroPartialId, pc.GidxPartialId)
	if err != nil {
		t.Fatalf("Error counting by partial comparison: %s\n", err.Error())
	}

	if num != int64(1) {
		t.Fatalf("Partial comparison count incorrect\n")
	}
}

func TestPartialComparisonServiceFindAll(t *testing.T) {
	setupPartialComparisonServiceTest()
	partialComparisonService := serviceFactory.MustPartialComparisonService()
	defer partialComparisonService.Close()

	pc := model.PartialComparison{
		MacroPartialId: macroPartial.Id,
		GidxPartialId:  gidxPartial.Id,
		Dist:           0.5,
	}

	err := partialComparisonService.Insert(&pc)
	if err != nil {
		t.Fatalf("Error inserting partial comparison: %s\n", err.Error())
	}

	pcs, err := partialComparisonService.FindAll("id DESC", 1000, 0, "macro_partial_id = ?", macroPartial.Id)
	if err != nil {
		t.Fatalf("Error finding all partial comparisons: %s\n", err.Error())
	}

	if pcs == nil {
		t.Fatalf("No partial comparison slice returned for FindAll\n")
	}

	if len(pcs) != 1 {
		t.Fatal("Inserted partial comparison not found by FindAll")
	}

	pc2 := pcs[0]

	if pc2.MacroPartialId != pc.MacroPartialId ||
		pc2.GidxPartialId != pc.GidxPartialId {
		t.Fatal("partial comparison macro id does not match")
	}
}

func TestPartialComparisonServiceFindOrCreate(t *testing.T) {
	setupPartialComparisonServiceTest()
	partialComparisonService := serviceFactory.MustPartialComparisonService()
	defer partialComparisonService.Close()

	partialComparison, err := partialComparisonService.FindOrCreate(&macroPartial, &gidxPartial)
	if err != nil {
		t.Fatalf("Failed to FindOrCreate partialComparison: %s\n", err.Error())
	}

	if partialComparison.MacroPartialId != macroPartial.Id {
		t.Fatalf("partialComparison.MacroPartialId was %d, expected %d\n", partialComparison.MacroPartialId, macroPartial.Id)
	}

	if partialComparison.GidxPartialId != gidxPartial.Id {
		t.Fatalf("partialComparison.GidxPartialId was %d, expected %d\n", partialComparison.GidxPartialId, gidxPartial.Id)
	}

	if partialComparison.Dist == 0.0 {
		t.Fatalf("partial comparison dist was 0.0")
	}
}

func TestPartialComparisonServiceCountMissing(t *testing.T) {
	setupPartialComparisonServiceTest()
	partialComparisonService := serviceFactory.MustPartialComparisonService()
	defer partialComparisonService.Close()

	num, err := partialComparisonService.CountMissing(&macro)
	if err != nil {
		t.Fatalf("Error counting missing partial comparisons: %s\n", err.Error())
	}

	if num != 10 {
		t.Fatalf("Expected 10 missing partial comparisons, got %d\n", num)
	}
}

func TestPartialComparisonServiceFindMissing(t *testing.T) {
	setupPartialComparisonServiceTest()
	partialComparisonService := serviceFactory.MustPartialComparisonService()
	defer partialComparisonService.Close()

	macroGidxViews, err := partialComparisonService.FindMissing(&macro, 1000)
	if err != nil {
		t.Fatalf("Error finding missing partial comparisons: %s\n", err.Error())
	}

	if len(macroGidxViews) != 10 {
		t.Fatalf("Expected 10 missing partial comparisons, got %d\n", len(macroGidxViews))
	}
}

func TestPartialComparisonServiceCreateFromView(t *testing.T) {
	setupPartialComparisonServiceTest()
	partialComparisonService := serviceFactory.MustPartialComparisonService()
	defer partialComparisonService.Close()

	macroGidxViews, err := partialComparisonService.FindMissing(&macro, 1000)
	if err != nil {
		t.Fatalf("Error finding missing partial comparisons: %s\n", err.Error())
	}

	if len(macroGidxViews) != 10 {
		t.Fatalf("Expected 10 missing partial comparisons, got %d\n", len(macroGidxViews))
	}

	view := macroGidxViews[0]
	pc, err := partialComparisonService.CreateFromView(view)
	if err != nil {
		t.Fatalf("Error creating partial comparison from view: %s\n", err.Error())
	}

	if pc == nil {
		t.Fatal("Partial comparison not created from view")
	}

	if pc.Id == int64(0) {
		t.Fatal("Partial comparison from view not given id")
	}

	if pc.Dist == 0.0 {
		t.Fatal("Partial comparison dist not calculated")
	}
}

func TestPartialComparisonServiceGetClosest(t *testing.T) {
	setupPartialComparisonServiceTest()
	partialComparisonService := serviceFactory.MustPartialComparisonService()
	defer partialComparisonService.Close()

	pc1 := model.PartialComparison{
		MacroPartialId: macroPartial.Id,
		GidxPartialId:  int64(1),
		Dist:           0.2,
	}
	err := partialComparisonService.Insert(&pc1)
	if err != nil {
		t.Fatalf("Error inserting partial comparison: %s\n", err.Error())
	}

	pc2 := model.PartialComparison{
		MacroPartialId: macroPartial.Id,
		GidxPartialId:  int64(2),
		Dist:           0.1,
	}
	err = partialComparisonService.Insert(&pc2)
	if err != nil {
		t.Fatalf("Error inserting partial comparison: %s\n", err.Error())
	}

	gidxPartialId, err := partialComparisonService.GetClosest(&macroPartial)
	if err != nil {
		t.Fatalf("Error getting closest partial comparison: %s\n", err.Error())
	}

	if gidxPartialId != int64(2) {
		t.Fatalf("Expected closest partial comparison to have gidx id 2, but got %s\n", gidxPartialId)
	}
}

func TestPartialComparisonServiceGetClosestMax(t *testing.T) {
	setupPartialComparisonServiceTest()
	partialComparisonService := serviceFactory.MustPartialComparisonService()
	mosaicService := serviceFactory.MustMosaicService()
	mosaicPartialService := serviceFactory.MustMosaicPartialService()
	defer partialComparisonService.Close()

	mosaic = model.Mosaic{
		MacroId: macro.Id,
	}
	err := mosaicService.Insert(&mosaic)
	if err != nil {
		t.Fatalf("Error inserting mosaic: %s\n", err.Error())
	}

	pc1 := model.PartialComparison{
		MacroPartialId: macroPartial.Id,
		GidxPartialId:  int64(1),
		Dist:           0.1,
	}
	err = partialComparisonService.Insert(&pc1)
	if err != nil {
		t.Fatalf("Error inserting partial comparison: %s\n", err.Error())
	}

	pc2 := model.PartialComparison{
		MacroPartialId: macroPartial.Id,
		GidxPartialId:  int64(2),
		Dist:           0.2,
	}
	err = partialComparisonService.Insert(&pc2)
	if err != nil {
		t.Fatalf("Error inserting partial comparison: %s\n", err.Error())
	}

	// create mosaic partial with gidx 1, which has the closer distance,
	// GetClosestMax should then skip gidx 1 due to the maxRepeats argument
	mp1 := model.MosaicPartial{
		MosaicId:       mosaic.Id,
		MacroPartialId: macroPartial.Id,
		GidxPartialId:  int64(1),
	}
	err = mosaicPartialService.Insert(&mp1)
	if err != nil {
		t.Fatalf("Error inserting mosaic partial: %s\n", err.Error())
	}

	gidxPartialId, err := partialComparisonService.GetClosestMax(&macroPartial, &mosaic, 1)
	if err != nil {
		t.Fatalf("Error getting closest partial comparison: %s\n", err.Error())
	}

	if gidxPartialId != int64(2) {
		t.Fatalf("Expected closest partial comparison to have gidx id 2, but got %d\n", gidxPartialId)
	}
}

func TestPartialComparisonServiceGetBestAvailable(t *testing.T) {
	setupPartialComparisonServiceTest()
	partialComparisonService := serviceFactory.MustPartialComparisonService()
	defer partialComparisonService.Close()

	pc1 := model.PartialComparison{
		MacroPartialId: macroPartial.Id,
		GidxPartialId:  int64(1),
		Dist:           0.2,
	}
	err := partialComparisonService.Insert(&pc1)
	if err != nil {
		t.Fatalf("Error inserting partial comparison: %s\n", err.Error())
	}

	pc2 := model.PartialComparison{
		MacroPartialId: macroPartial.Id,
		GidxPartialId:  int64(2),
		Dist:           0.1,
	}
	err = partialComparisonService.Insert(&pc2)
	if err != nil {
		t.Fatalf("Error inserting partial comparison: %s\n", err.Error())
	}

	partialComparison, err := partialComparisonService.GetBestAvailable(&mosaic)
	if err != nil {
		t.Fatalf("Error getting best available partial comparison: %s\n", err.Error())
	} else if partialComparison == nil {
		t.Fatal("Partial comparison not found")
	}

	if partialComparison.GidxPartialId != int64(2) {
		t.Fatalf("Expected best partial comparison to have gidx id 2, but got %s\n", partialComparison.GidxPartialId)
	}
}

func TestPartialComparisonServiceGetBestAvailableMax(t *testing.T) {
	setupPartialComparisonServiceTest()
	partialComparisonService := serviceFactory.MustPartialComparisonService()
	mosaicService := serviceFactory.MustMosaicService()
	mosaicPartialService := serviceFactory.MustMosaicPartialService()
	macroPartialService := serviceFactory.MustMacroPartialService()
	defer partialComparisonService.Close()

	mosaic = model.Mosaic{
		MacroId: macro.Id,
	}
	err := mosaicService.Insert(&mosaic)
	if err != nil {
		t.Fatalf("Error inserting mosaic: %s\n", err.Error())
	}

	macroPartials, err := macroPartialService.FindAll("macro_partials.id asc", 9, 0, "macro_id = ?", macro.Id)
	if err != nil {
		t.Fatalf("Error getting macro partials: %s\n", err.Error())
	}

	for i, mp := range macroPartials {
		var dist float64
		var gidxId int64
		switch i {
		default:
			dist = 0.9
			gidxId = int64(1)
		case 1:
			dist = 0.2
			gidxId = int64(2)
		case 0:
			dist = 0.1
			gidxId = int64(1)
		}
		pc := model.PartialComparison{
			MacroPartialId: mp.Id,
			GidxPartialId:  gidxId,
			Dist:           dist,
		}
		err = partialComparisonService.Insert(&pc)
		if err != nil {
			t.Fatalf("Error inserting partial comparison: %s\n", err.Error())
		}
	}

	// create mosaic partial with gidx 1, which has the closer distance,
	// GetBestAvailableMax should then skip gidx 1 due to the maxRepeats argument
	mp1 := model.MosaicPartial{
		MosaicId:       mosaic.Id,
		MacroPartialId: macroPartials[0].Id,
		GidxPartialId:  int64(1),
	}
	err = mosaicPartialService.Insert(&mp1)
	if err != nil {
		t.Fatalf("Error inserting mosaic partial: %s\n", err.Error())
	}

	partialComparison, err := partialComparisonService.GetBestAvailableMax(&mosaic, 1)
	if err != nil {
		t.Fatalf("Error getting best available partial comparison: %s\n", err.Error())
	} else if partialComparison == nil {
		t.Fatal("Partial comparison not found")
	}

	if partialComparison.GidxPartialId != int64(2) {
		t.Fatalf("Expected best partial comparison to have gidx id 2, but got %d\n", partialComparison.GidxPartialId)
	}
}
