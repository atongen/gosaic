package service

import (
	"testing"

	"gosaic/model"

	_ "github.com/mattn/go-sqlite3"
)

func setupCoverPartialServiceTest() (CoverPartialService, error) {
	dbMap, err := getTestDbMap()
	if err != nil {
		return nil, err
	}

	coverService, err := getTestCoverService(dbMap)
	if err != nil {
		return nil, err
	}

	aspectService, err := getTestAspectService(dbMap)
	if err != nil {
		return nil, err
	}

	coverPartialService, err := getTestCoverPartialService(dbMap)
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

	return coverPartialService, nil
}

func TestCoverPartialServiceInsert(t *testing.T) {
	coverPartialService, err := setupCoverPartialServiceTest()
	if err != nil {
		t.Fatalf("Unable to setup services: %s\n", err.Error())
	}
	defer coverPartialService.Close()

	p1 := model.CoverPartial{
		CoverId:  cover.Id,
		AspectId: aspect.Id,
		X1:       0,
		Y1:       0,
		X2:       2,
		Y2:       2,
	}

	err = coverPartialService.Insert(&p1)
	if err != nil {
		t.Fatalf("Error inserting cover partial: %s\n", err.Error())
	}

	if p1.Id == int64(0) {
		t.Fatalf("Inserted cover partial id not set")
	}

	p2, err := coverPartialService.Get(p1.Id)
	if err != nil {
		t.Fatalf("Error getting inserted cover partial: %s\n", err.Error())
	} else if p2 == nil {
		t.Fatalf("Cover partial not inserted\n")
	}

	if p1.Id != p2.Id ||
		p1.CoverId != p2.CoverId ||
		p1.AspectId != p2.AspectId ||
		p1.X1 != p2.X1 ||
		p1.Y1 != p2.Y1 ||
		p1.X2 != p2.X2 ||
		p1.Y2 != p2.Y2 {
		t.Fatalf("Inserted cover partial (%+v) does not match: %+v\n", p2, p1)
	}
}

func TestCoverPartialServiceBulkInsert(t *testing.T) {
	coverPartialService, err := setupCoverPartialServiceTest()
	if err != nil {
		t.Fatalf("Unable to setup services: %s\n", err.Error())
	}
	defer coverPartialService.Close()

	coverPartials := make([]*model.CoverPartial, 5)
	for i := 0; i < 5; i++ {
		coverPartials[i] = &model.CoverPartial{
			CoverId:  cover.Id,
			AspectId: aspect.Id,
			X1:       int64(i),
			Y1:       int64(i),
			X2:       int64(i + 1),
			Y2:       int64(i + 1),
		}
	}

	num, err := coverPartialService.BulkInsert(coverPartials)
	if err != nil {
		t.Fatalf("Error bulk inserting cover partials: %s\n", err.Error())
	}

	if num != 5 {
		t.Fatalf("Expected bulk insert result to be 5, but got %d\n", num)
	}

	icps, err := coverPartialService.FindAll(cover.Id, "id asc")
	if err != nil {
		t.Fatalf("Error finding bulk inserted cover partials: %s\n", err.Error())
	}

	if len(icps) != 5 {
		t.Fatalf("Expected 5 bulk inserted cover partials, got %d\n", len(icps))
	}
}

func TestCoverPartialServiceUpdate(t *testing.T) {
	coverPartialService, err := setupCoverPartialServiceTest()
	if err != nil {
		t.Fatalf("Unable to setup services: %s\n", err.Error())
	}
	defer coverPartialService.Close()

	p1 := model.CoverPartial{
		CoverId:  cover.Id,
		AspectId: aspect.Id,
		X1:       0,
		Y1:       0,
		X2:       2,
		Y2:       2,
	}

	err = coverPartialService.Insert(&p1)
	if err != nil {
		t.Fatalf("Error inserting cover partial: %s\n", err.Error())
	}

	p1.X1 = 1
	err = coverPartialService.Update(&p1)
	if err != nil {
		t.Fatalf("Error updating cover partial: %s\n", err.Error())
	}

	p2, err := coverPartialService.Get(p1.Id)
	if err != nil {
		t.Fatalf("Error getting inserted cover partial: %s\n", err.Error())
	}

	if p2.X1 != 1 {
		t.Fatalf("Error updating cover partial, expected x1 1, got x1 %d\n", p2.X1)
	}
}

func TestCoverPartialServiceDelete(t *testing.T) {
	coverPartialService, err := setupCoverPartialServiceTest()
	if err != nil {
		t.Fatalf("Unable to setup services: %s\n", err.Error())
	}
	defer coverPartialService.Close()

	p1 := model.CoverPartial{
		CoverId:  cover.Id,
		AspectId: aspect.Id,
		X1:       0,
		Y1:       0,
		X2:       2,
		Y2:       2,
	}

	err = coverPartialService.Insert(&p1)
	if err != nil {
		t.Fatalf("Error inserting cover partial: %s\n", err.Error())
	}

	err = coverPartialService.Delete(&p1)
	if err != nil {
		t.Fatalf("Error deleting cover partial: %s\n", err.Error())
	}

	p2, err := coverPartialService.Get(p1.Id)
	if err != nil {
		t.Fatalf("Error getting cover partial: %s\n", err.Error())
	} else if p2 != nil {
		t.Fatalf("Cover partial not deleted")
	}
}

func TestCoverPartialServiceFindAll(t *testing.T) {
	coverPartialService, err := setupCoverPartialServiceTest()
	if err != nil {
		t.Fatalf("Unable to setup services: %s\n", err.Error())
	}
	defer coverPartialService.Close()

	cps := []model.CoverPartial{
		model.CoverPartial{CoverId: cover.Id, AspectId: aspect.Id, X1: 0, Y1: 0, X2: 1, Y2: 1},
		model.CoverPartial{CoverId: cover.Id, AspectId: aspect.Id, X1: 0, Y1: 0, X2: 1, Y2: 1},
		model.CoverPartial{CoverId: cover.Id, AspectId: aspect.Id, X1: 0, Y1: 0, X2: 1, Y2: 1},
	}

	for _, cp := range cps {
		err = coverPartialService.Insert(&cp)
		if err != nil {
			t.Fatalf("Error inserting cover partial: %s\n", err.Error())
		}
	}

	cps2, err := coverPartialService.FindAll(cover.Id, "cover_partials.id ASC")
	if err != nil {
		t.Fatalf("Error finding cover partials: %s\n", err.Error())
	}

	if len(cps2) != 3 {
		t.Fatalf("Wanted 3 cover partials, got %d\n", len(cps))
	}
}
