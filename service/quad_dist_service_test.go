package service

import (
	"testing"

	"github.com/atongen/gosaic/model"

	_ "github.com/mattn/go-sqlite3"
)

func setupQuadDistServiceTest() {
	setTestServiceFactory()
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

	cover = model.Cover{AspectId: aspect.Id, Width: 1, Height: 1}
	err = coverService.Insert(&cover)
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

	coverPartials := make([]model.CoverPartial, 0)
	for i := 0; i < 3; i++ {
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
		if i == 0 {
			coverPartial = cp
		} else {
			coverPartials = append(coverPartials, cp)
		}
	}

	for i, cp := range coverPartials {
		mp, err := macroPartialService.FindOrCreate(&macro, &cp)
		if err != nil {
			panic(err)
		}
		if i == 0 {
			macroPartial = *mp
		}
	}
}

func TestQuadDistServiceInsert(t *testing.T) {
	setupQuadDistServiceTest()
	quadDistService := serviceFactory.MustQuadDistService()
	defer quadDistService.Close()

	pc := model.QuadDist{
		MacroPartialId: macroPartial.Id,
		Depth:          10,
		Area:           100,
		Dist:           0.5,
	}

	err := quadDistService.Insert(&pc)
	if err != nil {
		t.Fatalf("Error inserting quad dist: %s\n", err.Error())
	}

	if pc.Id == int64(0) {
		t.Fatalf("Inserted quad dist id not set")
	}

	pc2, err := quadDistService.Get(pc.Id)
	if err != nil {
		t.Fatalf("Error getting inserted quad dist: %s\n", err.Error())
	} else if pc2 == nil {
		t.Fatalf("quad dist not inserted\n")
	}

	if pc.Id != pc2.Id ||
		pc.MacroPartialId != pc2.MacroPartialId ||
		pc.Area != pc2.Area ||
		pc.Depth != pc2.Depth ||
		pc.Dist != pc2.Dist {
		t.Fatal("Inserted macro partial data does not match")
	}
}

func TestQuadDistServiceGetWorst(t *testing.T) {
	setupQuadDistServiceTest()
	quadDistService := serviceFactory.MustQuadDistService()
	defer quadDistService.Close()

	pc1 := model.QuadDist{
		MacroPartialId: int64(1),
		Depth:          2,
		Area:           10,
		Dist:           0.4,
	}

	err := quadDistService.Insert(&pc1)
	if err != nil {
		t.Fatalf("Error inserting quad dist: %s\n", err.Error())
	}

	pc2 := model.QuadDist{
		MacroPartialId: int64(2),
		Depth:          2,
		Area:           10,
		Dist:           0.6,
	}

	err = quadDistService.Insert(&pc2)
	if err != nil {
		t.Fatalf("Error inserting quad dist: %s\n", err.Error())
	}

	coverPartialQuadView, err := quadDistService.GetWorst(&macro, 100, 0)
	if err != nil {
		t.Fatalf("Error getting worst quad dist: %s\n", err.Error())
	} else if coverPartialQuadView == nil {
		t.Fatal("worst quad dist not found")
	}

	// id 3 corresponds to 2nd macro partial
	if coverPartialQuadView.CoverPartial.Id != int64(3) {
		t.Fatalf("Expected cover partial id 3 to be worst, got %d\n", coverPartialQuadView.CoverPartial.Id)
	}

}
