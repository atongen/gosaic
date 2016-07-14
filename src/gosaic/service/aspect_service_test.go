package service

import (
	"testing"

	"gosaic/model"

	_ "github.com/mattn/go-sqlite3"
)

var (
	testAspect1 = model.NewAspect(1, 1)
)

func setupAspectServiceTest() (AspectService, error) {
	dbMap, err := getTestDbMap()
	if err != nil {
		return nil, err
	}

	aspectService, err := getTestAspectService(dbMap)
	if err != nil {
		return nil, err
	}

	aspect := model.NewAspect(testAspect1.Columns, testAspect1.Rows)
	err = aspectService.Insert(aspect)
	if err != nil {
		return nil, err
	}
	testAspect1.Id = aspect.Id
	return aspectService, nil
}

func TestAspectServiceGet(t *testing.T) {
	aspectService, err := setupAspectServiceTest()
	if err != nil {
		t.Error("Unable to setup database.", err)
	}
	defer aspectService.DbMap().Db.Close()

	aspect, err := aspectService.Get(testAspect1.Id)
	if err != nil {
		t.Error("Error finding aspect by id", err)
	}

	if aspect.Id != testAspect1.Id ||
		aspect.Columns != testAspect1.Columns ||
		aspect.Rows != testAspect1.Rows {
		t.Error("Found aspect does not match data")
	}
}

func TestAspectServiceGetMissing(t *testing.T) {
	aspectService, err := setupAspectServiceTest()
	if err != nil {
		t.Error("Unable to setup database.", err)
	}
	defer aspectService.DbMap().Db.Close()

	aspect, err := aspectService.Get(1234)
	if err != nil {
		t.Error("Error finding aspect by id", err)
	}

	if aspect != nil {
		t.Error("Found non-existent aspect")
	}
}

func TestAspectServiceFindOrCreate(t *testing.T) {
	aspectService, err := setupAspectServiceTest()
	if err != nil {
		t.Error("Unable to setup database.", err)
	}
	defer aspectService.DbMap().Db.Close()

	n1, err := aspectService.Count()
	if err != nil {
		t.Error("Unable to count aspects")
	}

	a1, err := aspectService.FindOrCreate(100, 100)
	if err != nil {
		t.Error("Unable to find or create 100x100 aspect")
	}

	a2, err := aspectService.FindOrCreate(200, 200)
	if err != nil {
		t.Error("Unable to find or create 200x200 aspect")
	}

	a3, err := aspectService.FindOrCreate(300, 300)
	if err != nil {
		t.Error("Unable to find or create 300x300 aspect")
	}

	n2, err := aspectService.Count()
	if err != nil {
		t.Error("Unable to re-count aspects")
	}

	if n1 != n2 {
		t.Error("Created aspect when shouldn't have")
	}

	if a1.Id != a2.Id {
		t.Error("Aspects not equal")
	}

	if a2.Id != a3.Id {
		t.Error("Aspects not equal")
	}
}
