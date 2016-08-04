package service

import (
	"fmt"
	"testing"

	"gosaic/model"

	_ "github.com/mattn/go-sqlite3"
)

func setupMarcroPartialServiceTest() (MacroPartialService, error) {
	dbMap, err := getTestDbMap()
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

	aspect = model.Aspect{Columns: 1, Rows: 1}
	err = aspectService.Insert(&aspect)
	if err != nil {
		return nil, err
	}

	cover = model.Cover{Name: "test1", AspectId: aspect.Id, Width: 1, Height: 1}
	err = coverService.Insert(&cover)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	macro = model.Macro{
		AspectId:    aspect.Id,
		CoverId:     cover.Id,
		Path:        "/path/to/my/macro.jpg",
		Md5sum:      "d41d8cd98f00b204e9800998ecf8427e",
		Width:       1,
		Height:      1,
		Orientation: 1,
	}
	err = macroService.Insert(&macro)
	if err != nil {
		return nil, err
	}

	return macroPartialService, nil
}

func TestMacroPartialServiceInsert(t *testing.T) {
	macroPartialService, err := setupMarcroPartialServiceTest()
	if err != nil {
		t.Fatalf("Unable to setup database: %s\n", err.Error())
	}
	defer macroPartialService.DbMap().Db.Close()

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
	mp.EncodePixels()

	err = macroPartialService.Insert(&mp)
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

	fmt.Printf("mp2: %+s\n", mp2)
	fmt.Println(len(mp2.Pixels))

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
