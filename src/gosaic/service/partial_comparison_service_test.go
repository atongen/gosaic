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

func TestPartialComparisonServiceInsert(t *testing.T) {
	partialComparisonService, err := setupPartialComparisonServiceTest()
	if err != nil {
		t.Fatalf("Unable to setup database: %s\n", err.Error())
	}
	defer partialComparisonService.DbMap().Db.Close()

	pc := model.PartialComparison{
		MacroPartialId: macroPartial.Id,
		GidxPartialId:  gidxPartial.Id,
		Dist:           0.5,
	}

	err = partialComparisonService.Insert(&pc)
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

//func TestMacroPartialServiceUpdate(t *testing.T) {
//	macroPartialService, err := setupMacroPartialServiceTest()
//	if err != nil {
//		t.Fatalf("Unable to setup database: %s\n", err.Error())
//	}
//	defer macroPartialService.DbMap().Db.Close()
//
//	mp := model.MacroPartial{
//		MacroId:        macro.Id,
//		CoverPartialId: coverPartial.Id,
//		AspectId:       aspect.Id,
//		Pixels: []*model.Lab{
//			&model.Lab{
//				L:     0.4,
//				A:     0.5,
//				B:     0.6,
//				Alpha: 0.0,
//			},
//		},
//	}
//
//	err = macroPartialService.Insert(&mp)
//	if err != nil {
//		t.Fatalf("Error inserting macro partial: %s\n", err.Error())
//	}
//
//	mp.Pixels[0].L = 0.75
//	err = macroPartialService.Update(&mp)
//	if err != nil {
//		t.Fatalf("Error updating macro partial: %s\n", err.Error())
//	}
//
//	mp2, err := macroPartialService.Get(mp.Id)
//	if err != nil {
//		t.Fatalf("Error getting updated macro partial: %s\n", err.Error())
//	} else if mp2 == nil {
//		t.Fatalf("Macro partial not inserted\n")
//	}
//
//	if mp2.Pixels[0].L != 0.75 {
//		t.Fatal("Updated macro partial data does not match")
//	}
//}
//
//func TestMacroPartialServiceDelete(t *testing.T) {
//	macroPartialService, err := setupMacroPartialServiceTest()
//	if err != nil {
//		t.Fatalf("Unable to setup database: %s\n", err.Error())
//	}
//	defer macroPartialService.DbMap().Db.Close()
//
//	mp := model.MacroPartial{
//		MacroId:        macro.Id,
//		CoverPartialId: coverPartial.Id,
//		AspectId:       aspect.Id,
//		Pixels: []*model.Lab{
//			&model.Lab{
//				L:     0.4,
//				A:     0.5,
//				B:     0.6,
//				Alpha: 0.0,
//			},
//		},
//	}
//
//	err = macroPartialService.Insert(&mp)
//	if err != nil {
//		t.Fatalf("Error inserting macro partial: %s\n", err.Error())
//	}
//
//	err = macroPartialService.Delete(&mp)
//	if err != nil {
//		t.Fatalf("Error deleting macro partial: %s\n", err.Error())
//	}
//
//	mp2, err := macroPartialService.Get(mp.Id)
//	if err != nil {
//		t.Fatalf("Error getting deleted macro partial: %s\n", err.Error())
//	} else if mp2 != nil {
//		t.Fatalf("Macro partial not deleted\n")
//	}
//}
//
//func TestMacroPartialServiceGetOneBy(t *testing.T) {
//	macroPartialService, err := setupMacroPartialServiceTest()
//	if err != nil {
//		t.Fatalf("Unable to setup database: %s\n", err.Error())
//	}
//	defer macroPartialService.DbMap().Db.Close()
//
//	mp := model.MacroPartial{
//		MacroId:        macro.Id,
//		CoverPartialId: coverPartial.Id,
//		AspectId:       aspect.Id,
//		Pixels: []*model.Lab{
//			&model.Lab{
//				L:     0.4,
//				A:     0.5,
//				B:     0.6,
//				Alpha: 0.0,
//			},
//		},
//	}
//
//	err = macroPartialService.Insert(&mp)
//	if err != nil {
//		t.Fatalf("Error inserting macro partial: %s\n", err.Error())
//	}
//
//	mp2, err := macroPartialService.GetOneBy("macro_id", mp.MacroId)
//	if err != nil {
//		t.Fatalf("Error getting one by macro partial: %s\n", err.Error())
//	} else if mp2 == nil {
//		t.Fatalf("Macro partial not found by\n")
//	}
//
//	if mp2.MacroId != mp.MacroId {
//		t.Fatal("Macro partial macro id does not match")
//	}
//
//	if len(mp2.Pixels) != 1 {
//		t.Fatalf("Expected 1 macro partial pixel, got %d\n", len(mp2.Pixels))
//	}
//
//	plab := mp2.Pixels[0]
//
//	if plab.L != 0.4 &&
//		plab.A != 0.5 &&
//		plab.B != 0.6 &&
//		plab.Alpha != 0.0 {
//		t.Fatal("Macro partial pixel data is not correct")
//	}
//}
//
//func TestMacroPartialServiceExistsBy(t *testing.T) {
//	macroPartialService, err := setupMacroPartialServiceTest()
//	if err != nil {
//		t.Fatalf("Unable to setup database: %s\n", err.Error())
//	}
//	defer macroPartialService.DbMap().Db.Close()
//
//	mp := model.MacroPartial{
//		MacroId:        macro.Id,
//		CoverPartialId: coverPartial.Id,
//		AspectId:       aspect.Id,
//		Pixels: []*model.Lab{
//			&model.Lab{
//				L:     0.4,
//				A:     0.5,
//				B:     0.6,
//				Alpha: 0.0,
//			},
//		},
//	}
//
//	err = macroPartialService.Insert(&mp)
//	if err != nil {
//		t.Fatalf("Error inserting macro partial: %s\n", err.Error())
//	}
//
//	found, err := macroPartialService.ExistsBy("macro_id", mp.MacroId)
//	if err != nil {
//		t.Fatalf("Error getting one by macro partial: %s\n", err.Error())
//	}
//
//	if !found {
//		t.Fatalf("Macro partial not exists by\n")
//	}
//}
//
//func TestMacroPartialServiceCount(t *testing.T) {
//	macroPartialService, err := setupMacroPartialServiceTest()
//	if err != nil {
//		t.Fatalf("Unable to setup database: %s\n", err.Error())
//	}
//	defer macroPartialService.DbMap().Db.Close()
//
//	mp := model.MacroPartial{
//		MacroId:        macro.Id,
//		CoverPartialId: coverPartial.Id,
//		AspectId:       aspect.Id,
//		Pixels: []*model.Lab{
//			&model.Lab{
//				L:     0.4,
//				A:     0.5,
//				B:     0.6,
//				Alpha: 0.0,
//			},
//		},
//	}
//
//	err = macroPartialService.Insert(&mp)
//	if err != nil {
//		t.Fatalf("Error inserting macro partial: %s\n", err.Error())
//	}
//
//	num, err := macroPartialService.Count()
//	if err != nil {
//		t.Fatalf("Error counting macro partial: %s\n", err.Error())
//	}
//
//	if num != int64(1) {
//		t.Fatalf("Macro partial count incorrect\n")
//	}
//}
//
//func TestMacroPartialServiceCountBy(t *testing.T) {
//	macroPartialService, err := setupMacroPartialServiceTest()
//	if err != nil {
//		t.Fatalf("Unable to setup database: %s\n", err.Error())
//	}
//	defer macroPartialService.DbMap().Db.Close()
//
//	mp := model.MacroPartial{
//		MacroId:        macro.Id,
//		CoverPartialId: coverPartial.Id,
//		AspectId:       aspect.Id,
//		Pixels: []*model.Lab{
//			&model.Lab{
//				L:     0.4,
//				A:     0.5,
//				B:     0.6,
//				Alpha: 0.0,
//			},
//		},
//	}
//
//	err = macroPartialService.Insert(&mp)
//	if err != nil {
//		t.Fatalf("Error inserting macro partial: %s\n", err.Error())
//	}
//
//	num, err := macroPartialService.CountBy("macro_id", macro.Id)
//	if err != nil {
//		t.Fatalf("Error counting macro partial: %s\n", err.Error())
//	}
//
//	if num != int64(1) {
//		t.Fatalf("Macro partial count incorrect\n")
//	}
//}
//
//func TestMacroPartialServiceFindAll(t *testing.T) {
//	macroPartialService, err := setupMacroPartialServiceTest()
//	if err != nil {
//		t.Fatalf("Unable to setup database: %s\n", err.Error())
//	}
//	defer macroPartialService.DbMap().Db.Close()
//
//	mp := model.MacroPartial{
//		MacroId:        macro.Id,
//		CoverPartialId: coverPartial.Id,
//		AspectId:       aspect.Id,
//		Pixels: []*model.Lab{
//			&model.Lab{
//				L:     0.4,
//				A:     0.5,
//				B:     0.6,
//				Alpha: 0.0,
//			},
//		},
//	}
//
//	err = macroPartialService.Insert(&mp)
//	if err != nil {
//		t.Fatalf("Error inserting macro partial: %s\n", err.Error())
//	}
//
//	mps, err := macroPartialService.FindAll("id DESC", 1000, 0, "cover_partial_id = ?", coverPartial.Id)
//	if err != nil {
//		t.Fatalf("Error finding all macro partials: %s\n", err.Error())
//	}
//
//	if mps == nil {
//		t.Fatalf("No macro partial slice returned for FindAll\n")
//	}
//
//	if len(mps) != 1 {
//		t.Fatal("Inserted macro partial not found by FindAll")
//	}
//
//	mp2 := mps[0]
//
//	if mp2.MacroId != mp.MacroId {
//		t.Fatal("Macro partial macro id does not match")
//	}
//
//	if len(mp2.Pixels) != 1 {
//		t.Fatalf("Expected 1 macro partial pixel, got %d\n", len(mp2.Pixels))
//	}
//
//	plab := mp2.Pixels[0]
//
//	if plab.L != 0.4 &&
//		plab.A != 0.5 &&
//		plab.B != 0.6 &&
//		plab.Alpha != 0.0 {
//		t.Fatal("Macro partial pixel data is not correct")
//	}
//}
//
//func TestMacroPartialServiceFindOrCreate(t *testing.T) {
//	macroPartialService, err := setupMacroPartialServiceTest()
//	if err != nil {
//		t.Fatalf("Unable to setup database: %s\n", err.Error())
//	}
//	defer macroPartialService.DbMap().Db.Close()
//
//	macroPartial, err := macroPartialService.FindOrCreate(&macro, &coverPartial)
//	if err != nil {
//		t.Fatalf("Failed to FindOrCreate macroPartial: %s\n", err.Error())
//	}
//
//	if macroPartial.MacroId != macro.Id {
//		t.Errorf("macroPartial.MacroId was %d, expected %d\n", macroPartial.MacroId, macro.Id)
//	}
//
//	if macroPartial.CoverPartialId != coverPartial.Id {
//		t.Errorf("macroPartial.CoverPartialId was %d, expected %d\n", macroPartial.CoverPartialId, coverPartial.Id)
//	}
//
//	if macroPartial.AspectId != aspect.Id {
//		t.Errorf("macroPartial.AspectId was %d, expected %d\n", macroPartial.AspectId, aspect.Id)
//	}
//
//	if len(macroPartial.Data) == 0 {
//		t.Error("macroPartial.Data was empty")
//	}
//
//	numPixels := len(macroPartial.Pixels)
//	if numPixels != 100 {
//		t.Errorf("macroPartial.Pixels len was %d, expected %d\n", numPixels, 100)
//	}
//
//	for i, pix := range macroPartial.Pixels {
//		if pix.L == 0.0 && pix.A == 0.0 && pix.B == 0.0 && pix.Alpha == 0.0 {
//			t.Errorf("pixel %d was empty\n", i)
//		}
//	}
//}
//
//

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
