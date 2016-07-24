package service

import (
	"testing"

	"gosaic/model"

	_ "github.com/mattn/go-sqlite3"
)

func setupCoverServiceTest() (CoverService, error) {
	dbMap, err := getTestDbMap()
	if err != nil {
		return nil, err
	}

	coverService, err := getTestCoverService(dbMap)
	if err != nil {
		return nil, err
	}

	return coverService, nil
}

func TestCoverServiceInsert(t *testing.T) {
	coverService, err := setupCoverServiceTest()
	if err != nil {
		t.Fatalf("Unable to setup database: %s\n", err.Error())
	}
	defer coverService.DbMap().Db.Close()

	c1 := model.Cover{
		Name:   "test1",
		Type:   "test",
		Width:  600,
		Height: 400,
	}

	err = coverService.Insert(&c1)
	if err != nil {
		t.Fatalf("Error inserting cover: %s\n", err.Error())
	}

	if c1.Id == int64(0) {
		t.Fatalf("Inserted cover id not set")
	}

	c2, err := coverService.Get(c1.Id)
	if err != nil {
		t.Fatalf("Error getting inserted cover: %s\n", err.Error())
	} else if c2 == nil {
		t.Fatalf("Cover not inserted\n")
	}

	if c1.Id != c2.Id ||
		c1.Name != c2.Name ||
		c1.Type != c2.Type ||
		c1.Width != c2.Width ||
		c1.Height != c2.Height {
		t.Fatalf("Inserted cover (%+v) does not match: %+v\n", c2, c1)
	}
}

func TestCoverServiceUpdate(t *testing.T) {
	coverService, err := setupCoverServiceTest()
	if err != nil {
		t.Fatalf("Unable to setup database: %s\n", err.Error())
	}
	defer coverService.DbMap().Db.Close()

	c1 := model.Cover{
		Name:   "test1",
		Type:   "test",
		Width:  600,
		Height: 400,
	}

	err = coverService.Insert(&c1)
	if err != nil {
		t.Fatalf("Error inserting cover: %s\n", err.Error())
	}

	c1.Width = 800
	err = coverService.Update(&c1)
	if err != nil {
		t.Fatalf("Error updating cover: %s\n", err.Error())
	}

	c2, err := coverService.Get(c1.Id)
	if err != nil {
		t.Fatalf("Error getting inserted cover: %s\n", err.Error())
	}

	if c2.Width != 800 {
		t.Fatalf("Error updating cover, expected width 800, got width %s\n", c2.Width)
	}
}

func TestCoverServiceDelete(t *testing.T) {
	coverService, err := setupCoverServiceTest()
	if err != nil {
		t.Fatalf("Unable to setup database: %s\n", err.Error())
	}
	defer coverService.DbMap().Db.Close()

	c1 := model.Cover{
		Name:   "test1",
		Type:   "test",
		Width:  600,
		Height: 400,
	}

	err = coverService.Insert(&c1)
	if err != nil {
		t.Fatalf("Error inserting cover: %s\n", err.Error())
	}

	err = coverService.Delete(&c1)
	if err != nil {
		t.Fatalf("Error deleting cover: %s\n", err.Error())
	}

	c2, err := coverService.Get(c1.Id)
	if err != nil {
		t.Fatalf("Error getting cover: %s\n", err.Error())
	} else if c2 != nil {
		t.Fatalf("Cover not deleted")
	}
}

func TestCoverServiceGetOneBy(t *testing.T) {
	coverService, err := setupCoverServiceTest()
	if err != nil {
		t.Fatalf("Unable to setup database: %s\n", err.Error())
	}
	defer coverService.DbMap().Db.Close()

	c1 := model.Cover{
		Name:   "test1",
		Type:   "test",
		Width:  600,
		Height: 400,
	}

	err = coverService.Insert(&c1)
	if err != nil {
		t.Fatalf("Error inserting cover: %s\n", err.Error())
	}

	c2, err := coverService.GetOneBy("name", "test1")
	if err != nil {
		t.Fatalf("Error getting inserted cover: %s\n", err.Error())
	} else if c2 == nil {
		t.Fatalf("Cover not inserted\n")
	}

	if c1.Id != c2.Id ||
		c1.Name != c2.Name ||
		c1.Type != c2.Type ||
		c1.Width != c2.Width ||
		c1.Height != c2.Height {
		t.Fatalf("Inserted cover (%+v) does not match: %+v\n", c2, c1)
	}
}

func TestCoverServiceFindAll(t *testing.T) {
	coverService, err := setupCoverServiceTest()
	if err != nil {
		t.Fatalf("Unable to setup database: %s\n", err.Error())
	}
	defer coverService.DbMap().Db.Close()

	covers := []model.Cover{
		model.Cover{Name: "test1", Type: "test", Width: 600, Height: 400},
		model.Cover{Name: "test2", Type: "test", Width: 600, Height: 400},
		model.Cover{Name: "test3", Type: "test", Width: 600, Height: 400},
	}

	for _, cover := range covers {
		err = coverService.Insert(&cover)
		if err != nil {
			t.Fatalf("Error inserting cover: %s\n", err.Error())
		}
	}

	covers2, err := coverService.FindAll("covers.name ASC")
	if err != nil {
		t.Fatalf("Error finding covers: %s\n", err.Error())
	}

	if len(covers2) != 3 {
		t.Fatalf("Wanted 3 covers, got %d\n", len(covers2))
	}
}
