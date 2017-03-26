package service

import (
	"testing"

	"github.com/atongen/gosaic/model"

	_ "github.com/mattn/go-sqlite3"
)

func setupMacroPartialServiceTest() {
	setTestServiceFactory()
	aspectService := serviceFactory.MustAspectService()
	coverService := serviceFactory.MustCoverService()
	coverPartialService := serviceFactory.MustCoverPartialService()
	macroService := serviceFactory.MustMacroService()

	aspect = model.Aspect{Columns: 87, Rows: 128}
	err := aspectService.Insert(&aspect)
	if err != nil {
		panic(err)
	}

	cover = model.Cover{AspectId: aspect.Id, Width: 1, Height: 1}
	err = coverService.Insert(&cover)
	if err != nil {
		panic(err)
	}

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
		panic(err)
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
}

func TestMacroPartialServiceInsert(t *testing.T) {
	setupMacroPartialServiceTest()
	macroPartialService := serviceFactory.MustMacroPartialService()
	defer macroPartialService.Close()

	mp := model.MacroPartial{
		MacroId:        macro.Id,
		CoverPartialId: coverPartial.Id,
		AspectId:       aspect.Id,
		Pixels: []*model.Lab{
			&model.Lab{
				L:     0.4,
				A:     0.5,
				B:     0.6,
				Alpha: 0.0,
			},
		},
	}

	err := macroPartialService.Insert(&mp)
	if err != nil {
		t.Fatalf("Error inserting macro partial: %s\n", err.Error())
	}

	if mp.Id == int64(0) {
		t.Fatalf("Inserted macro partial id not set")
	}

	mp2, err := macroPartialService.Get(mp.Id)
	if err != nil {
		t.Fatalf("Error getting inserted macro partial: %s\n", err.Error())
	} else if mp2 == nil {
		t.Fatalf("Macro partial not inserted\n")
	}

	if mp.Id != mp2.Id ||
		mp.MacroId != mp2.MacroId ||
		mp.AspectId != mp2.AspectId {
		t.Fatal("Inserted macro partial data does not match")
	}

	if len(mp2.Pixels) != 1 {
		t.Fatal("Macro partial pixels not serialized correctly")
	}

	plab := mp2.Pixels[0]

	if plab.L != 0.4 &&
		plab.A != 0.5 &&
		plab.B != 0.6 &&
		plab.Alpha != 0.0 {
		t.Fatal("Macro partial pixel data is not correct")
	}
}

func TestMacroPartialServiceUpdate(t *testing.T) {
	setupMacroPartialServiceTest()
	macroPartialService := serviceFactory.MustMacroPartialService()
	defer macroPartialService.Close()

	mp := model.MacroPartial{
		MacroId:        macro.Id,
		CoverPartialId: coverPartial.Id,
		AspectId:       aspect.Id,
		Pixels: []*model.Lab{
			&model.Lab{
				L:     0.4,
				A:     0.5,
				B:     0.6,
				Alpha: 0.0,
			},
		},
	}

	err := macroPartialService.Insert(&mp)
	if err != nil {
		t.Fatalf("Error inserting macro partial: %s\n", err.Error())
	}

	mp.Pixels[0].L = 0.75
	err = macroPartialService.Update(&mp)
	if err != nil {
		t.Fatalf("Error updating macro partial: %s\n", err.Error())
	}

	mp2, err := macroPartialService.Get(mp.Id)
	if err != nil {
		t.Fatalf("Error getting updated macro partial: %s\n", err.Error())
	} else if mp2 == nil {
		t.Fatalf("Macro partial not inserted\n")
	}

	if mp2.Pixels[0].L != 0.75 {
		t.Fatal("Updated macro partial data does not match")
	}
}

func TestMacroPartialServiceDelete(t *testing.T) {
	setupMacroPartialServiceTest()
	macroPartialService := serviceFactory.MustMacroPartialService()
	defer macroPartialService.Close()

	mp := model.MacroPartial{
		MacroId:        macro.Id,
		CoverPartialId: coverPartial.Id,
		AspectId:       aspect.Id,
		Pixels: []*model.Lab{
			&model.Lab{
				L:     0.4,
				A:     0.5,
				B:     0.6,
				Alpha: 0.0,
			},
		},
	}

	err := macroPartialService.Insert(&mp)
	if err != nil {
		t.Fatalf("Error inserting macro partial: %s\n", err.Error())
	}

	err = macroPartialService.Delete(&mp)
	if err != nil {
		t.Fatalf("Error deleting macro partial: %s\n", err.Error())
	}

	mp2, err := macroPartialService.Get(mp.Id)
	if err != nil {
		t.Fatalf("Error getting deleted macro partial: %s\n", err.Error())
	} else if mp2 != nil {
		t.Fatalf("Macro partial not deleted\n")
	}
}

func TestMacroPartialServiceGetOneBy(t *testing.T) {
	setupMacroPartialServiceTest()
	macroPartialService := serviceFactory.MustMacroPartialService()
	defer macroPartialService.Close()

	mp := model.MacroPartial{
		MacroId:        macro.Id,
		CoverPartialId: coverPartial.Id,
		AspectId:       aspect.Id,
		Pixels: []*model.Lab{
			&model.Lab{
				L:     0.4,
				A:     0.5,
				B:     0.6,
				Alpha: 0.0,
			},
		},
	}

	err := macroPartialService.Insert(&mp)
	if err != nil {
		t.Fatalf("Error inserting macro partial: %s\n", err.Error())
	}

	mp2, err := macroPartialService.GetOneBy("macro_id", mp.MacroId)
	if err != nil {
		t.Fatalf("Error getting one by macro partial: %s\n", err.Error())
	} else if mp2 == nil {
		t.Fatalf("Macro partial not found by\n")
	}

	if mp2.MacroId != mp.MacroId {
		t.Fatal("Macro partial macro id does not match")
	}

	if len(mp2.Pixels) != 1 {
		t.Fatalf("Expected 1 macro partial pixel, got %d\n", len(mp2.Pixels))
	}

	plab := mp2.Pixels[0]

	if plab.L != 0.4 &&
		plab.A != 0.5 &&
		plab.B != 0.6 &&
		plab.Alpha != 0.0 {
		t.Fatal("Macro partial pixel data is not correct")
	}
}

func TestMacroPartialServiceExistsBy(t *testing.T) {
	setupMacroPartialServiceTest()
	macroPartialService := serviceFactory.MustMacroPartialService()
	defer macroPartialService.Close()

	mp := model.MacroPartial{
		MacroId:        macro.Id,
		CoverPartialId: coverPartial.Id,
		AspectId:       aspect.Id,
		Pixels: []*model.Lab{
			&model.Lab{
				L:     0.4,
				A:     0.5,
				B:     0.6,
				Alpha: 0.0,
			},
		},
	}

	err := macroPartialService.Insert(&mp)
	if err != nil {
		t.Fatalf("Error inserting macro partial: %s\n", err.Error())
	}

	found, err := macroPartialService.ExistsBy("macro_id", mp.MacroId)
	if err != nil {
		t.Fatalf("Error getting one by macro partial: %s\n", err.Error())
	}

	if !found {
		t.Fatalf("Macro partial not exists by\n")
	}
}

func TestMacroPartialServiceCount(t *testing.T) {
	setupMacroPartialServiceTest()
	macroPartialService := serviceFactory.MustMacroPartialService()
	defer macroPartialService.Close()

	mp := model.MacroPartial{
		MacroId:        macro.Id,
		CoverPartialId: coverPartial.Id,
		AspectId:       aspect.Id,
		Pixels: []*model.Lab{
			&model.Lab{
				L:     0.4,
				A:     0.5,
				B:     0.6,
				Alpha: 0.0,
			},
		},
	}

	err := macroPartialService.Insert(&mp)
	if err != nil {
		t.Fatalf("Error inserting macro partial: %s\n", err.Error())
	}

	num, err := macroPartialService.Count(&macro)
	if err != nil {
		t.Fatalf("Error counting macro partial: %s\n", err.Error())
	}

	if num != int64(1) {
		t.Fatalf("Macro partial count incorrect\n")
	}
}

func TestMacroPartialServiceFindAll(t *testing.T) {
	setupMacroPartialServiceTest()
	macroPartialService := serviceFactory.MustMacroPartialService()
	defer macroPartialService.Close()

	mp := model.MacroPartial{
		MacroId:        macro.Id,
		CoverPartialId: coverPartial.Id,
		AspectId:       aspect.Id,
		Pixels: []*model.Lab{
			&model.Lab{
				L:     0.4,
				A:     0.5,
				B:     0.6,
				Alpha: 0.0,
			},
		},
	}

	err := macroPartialService.Insert(&mp)
	if err != nil {
		t.Fatalf("Error inserting macro partial: %s\n", err.Error())
	}

	mps, err := macroPartialService.FindAll("id DESC", 1000, 0, "cover_partial_id = ?", coverPartial.Id)
	if err != nil {
		t.Fatalf("Error finding all macro partials: %s\n", err.Error())
	}

	if mps == nil {
		t.Fatalf("No macro partial slice returned for FindAll\n")
	}

	if len(mps) != 1 {
		t.Fatal("Inserted macro partial not found by FindAll")
	}

	mp2 := mps[0]

	if mp2.MacroId != mp.MacroId {
		t.Fatal("Macro partial macro id does not match")
	}

	if len(mp2.Pixels) != 1 {
		t.Fatalf("Expected 1 macro partial pixel, got %d\n", len(mp2.Pixels))
	}

	plab := mp2.Pixels[0]

	if plab.L != 0.4 &&
		plab.A != 0.5 &&
		plab.B != 0.6 &&
		plab.Alpha != 0.0 {
		t.Fatal("Macro partial pixel data is not correct")
	}
}

func TestMacroPartialServiceFindOrCreate(t *testing.T) {
	setupMacroPartialServiceTest()
	macroPartialService := serviceFactory.MustMacroPartialService()
	defer macroPartialService.Close()

	macroPartial, err := macroPartialService.FindOrCreate(&macro, &coverPartial)
	if err != nil {
		t.Fatalf("Failed to FindOrCreate macroPartial: %s\n", err.Error())
	}

	if macroPartial.MacroId != macro.Id {
		t.Errorf("macroPartial.MacroId was %d, expected %d\n", macroPartial.MacroId, macro.Id)
	}

	if macroPartial.CoverPartialId != coverPartial.Id {
		t.Errorf("macroPartial.CoverPartialId was %d, expected %d\n", macroPartial.CoverPartialId, coverPartial.Id)
	}

	if macroPartial.AspectId != aspect.Id {
		t.Errorf("macroPartial.AspectId was %d, expected %d\n", macroPartial.AspectId, aspect.Id)
	}

	if len(macroPartial.Data) == 0 {
		t.Error("macroPartial.Data was empty")
	}

	numPixels := len(macroPartial.Pixels)
	if numPixels != 100 {
		t.Errorf("macroPartial.Pixels len was %d, expected %d\n", numPixels, 100)
	}

	for i, pix := range macroPartial.Pixels {
		if pix.L == 0.0 && pix.A == 0.0 && pix.B == 0.0 && pix.Alpha == 0.0 {
			t.Errorf("pixel %d was empty\n", i)
		}
	}
}

func TestMacroPartialServiceCountMissing(t *testing.T) {
	setupMacroPartialServiceTest()
	macroPartialService := serviceFactory.MustMacroPartialService()
	defer macroPartialService.Close()

	num, err := macroPartialService.CountMissing(&macro)
	if err != nil {
		t.Fatalf("Error counting missing macro partials: %s\n", err.Error())
	}

	if num != 1 {
		t.Fatalf("Expected 1 missing macro partial, but got %d\n", num)
	}

	_, err = macroPartialService.Create(&macro, &coverPartial)
	if err != nil {
		t.Fatalf("Error creating missing macro partial: %s\n", err.Error())
	}

	num, err = macroPartialService.CountMissing(&macro)
	if err != nil {
		t.Fatalf("Error counting missing macro partials: %s\n", err.Error())
	}

	if num != 0 {
		t.Fatalf("Expected 0 missing macro partials, but got %d\n", num)
	}
}

func TestMacroPartialServiceFindMissing(t *testing.T) {
	setupMacroPartialServiceTest()
	macroPartialService := serviceFactory.MustMacroPartialService()
	defer macroPartialService.Close()

	coverPartials, err := macroPartialService.FindMissing(&macro, "id asc", 1000, 0)
	if err != nil {
		t.Fatalf("Error finding missing macro partials: %s\n", err.Error())
	}

	if len(coverPartials) != 1 {
		t.Fatalf("Expected 1 missing macro partial, but got %d\n", len(coverPartials))
	}

	mcp := coverPartials[0]

	if mcp.CoverId != coverPartial.CoverId ||
		mcp.AspectId != coverPartial.AspectId ||
		mcp.X1 != coverPartial.X1 ||
		mcp.Y1 != coverPartial.Y1 ||
		mcp.X2 != coverPartial.X2 ||
		mcp.Y2 != coverPartial.Y2 {
		t.Fatal("Missing macro partial does not match cover partial")
	}

	_, err = macroPartialService.Create(&macro, mcp)
	if err != nil {
		t.Fatalf("Error creating missing macro partial: %s\n", err.Error())
	}

	coverPartials, err = macroPartialService.FindMissing(&macro, "id asc", 1000, 0)
	if err != nil {
		t.Fatalf("Error finding missing macro partials: %s\n", err.Error())
	}

	if len(coverPartials) != 0 {
		t.Fatalf("Expected 0 missing macro partials, but got %d\n", len(coverPartials))
	}
}

func TestMacroPartialServiceAspectIds(t *testing.T) {
	setupMacroPartialServiceTest()
	macroPartialService := serviceFactory.MustMacroPartialService()
	defer macroPartialService.Close()

	_, err := macroPartialService.FindOrCreate(&macro, &coverPartial)
	if err != nil {
		t.Fatalf("Failed to FindOrCreate macroPartial: %s\n", err.Error())
	}

	aspectIds, err := macroPartialService.AspectIds(macro.AspectId)
	if err != nil {
		t.Fatalf("Failed to get aspect ids for macro partials: %s\n", err.Error())
	}

	macroPartials, err := macroPartialService.FindAll("id asc", 100, 0, "macro_id = ?", macro.Id)
	if err != nil {
		t.Fatalf("Failed to get all macro partials for macro: %s\n", err.Error())
	}

	if len(macroPartials) == 0 {
		t.Fatal("No macro partials found for macro")
	}

	for _, macroPartial := range macroPartials {
		found := false
		for _, aspectId := range aspectIds {
			if macroPartial.AspectId == aspectId {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("Expected to find aspect id %d in list of macro partials", macroPartial.AspectId)
		}
	}
}
