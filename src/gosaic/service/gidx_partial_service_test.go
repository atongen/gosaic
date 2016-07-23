package service

import (
	"testing"

	"gosaic/model"

	_ "github.com/mattn/go-sqlite3"
)

func TestGidxPartialServiceFindOrCreate(t *testing.T) {
	dbMap, err := getTestDbMap()
	if err != nil {
		t.Fatalf("Unable to get test dbmap: %s\n", err.Error())
	}
	defer dbMap.Db.Close()

	aspectService, err := getTestAspectService(dbMap)
	if err != nil {
		t.Fatalf("Unable to get aspect service: %s\n", err.Error())
	}

	gidxService, err := getTestGidxService(dbMap)
	if err != nil {
		t.Fatalf("Unable to get gidx service: %s\n", err.Error())
	}

	gidxPartialService, err := getTestGidxPartialService(dbMap)
	if err != nil {
		t.Fatalf("Unable to get gidx partial service: %s\n", err.Error())
	}

	aspect := model.Aspect{Columns: 87, Rows: 128}
	err = aspectService.Insert(&aspect)
	if err != nil {
		t.Fatalf("Unable to insert test aspect: %s\n", err.Error())
	}

	gidx := model.Gidx{
		AspectId:    aspect.Id,
		Path:        "testdata/matterhorn.jpg",
		Md5sum:      "fcaadee574094a3ae04c6badbbb9ee5e",
		Width:       uint(696),
		Height:      uint(1024),
		Orientation: 1,
	}
	err = gidxService.Insert(&gidx)
	if err != nil {
		t.Fatalf("Unable to insert test aspect: %s\n", err.Error())
	}

	gidxPartial, err := gidxPartialService.FindOrCreate(&gidx, &aspect)
	if err != nil {
		t.Fatalf("Failed to FindOrCreate gidxPartial: %s\n", err.Error())
	}

	if gidxPartial.GidxId != gidx.Id {
		t.Errorf("gidxPartial.GidxId was %d, expected %d\n", gidxPartial.GidxId, gidx.Id)
	}

	if gidxPartial.AspectId != aspect.Id {
		t.Errorf("gidxPartial.AspectId was %d, expected %d\n", gidxPartial.AspectId, aspect.Id)
	}

	if len(gidxPartial.Data) == 0 {
		t.Error("gidxPartial.Data was empty")
	}

	numPixels := len(gidxPartial.Pixels)
	if numPixels != 100 {
		t.Errorf("gidxPartial.Pixels len was %d, expected %d\n", numPixels, 100)
	}

	for i, pix := range gidxPartial.Pixels {
		if pix.L == 0.0 && pix.A == 0.0 && pix.B == 0.0 {
			t.Errorf("pixel %d was empty\n", i)
		}
	}
}

func TestGidxPartialServiceFindMissing(t *testing.T) {
	dbMap, err := getTestDbMap()
	if err != nil {
		t.Fatalf("Unable to get test dbmap: %s\n", err.Error())
	}
	defer dbMap.Db.Close()

	aspectService, err := getTestAspectService(dbMap)
	if err != nil {
		t.Fatalf("Unable to get aspect service: %s\n", err.Error())
	}

	gidxService, err := getTestGidxService(dbMap)
	if err != nil {
		t.Fatalf("Unable to get gidx service: %s\n", err.Error())
	}

	gidxPartialService, err := getTestGidxPartialService(dbMap)
	if err != nil {
		t.Fatalf("Unable to get gidx partial service: %s\n", err.Error())
	}

	aspect := model.Aspect{Columns: 87, Rows: 128}
	err = aspectService.Insert(&aspect)
	if err != nil {
		t.Fatalf("Unable to insert test aspect: %s\n", err.Error())
	}

	gidx := model.Gidx{
		AspectId:    aspect.Id,
		Path:        "testdata/matterhorn.jpg",
		Md5sum:      "fcaadee574094a3ae04c6badbbb9ee5e",
		Width:       uint(696),
		Height:      uint(1024),
		Orientation: 1,
	}
	err = gidxService.Insert(&gidx)
	if err != nil {
		t.Fatalf("Unable to insert test aspect: %s\n", err.Error())
	}

	gidxs, err := gidxPartialService.FindMissing(&aspect, "gidx.id", 100, 0)
	if err != nil {
		t.Fatalf("Failed to FindMissing gidxPartial: %s\n", err.Error())
	}

	if len(gidxs) != 1 {
		t.Fatalf("Expected 1 Missing gidxPartial, got %d\n", len(gidxs))
	}

	if gidxs[0].Id != gidx.Id {
		t.Errorf("Expected missing gidx id %d, got %d\n", gidx.Id, gidxs[0].Id)
	}
}
