package service

import (
	"testing"

	"github.com/atongen/gosaic/model"

	_ "github.com/mattn/go-sqlite3"
)

func setupCoverServiceTest() {
	setTestServiceFactory()
	aspectService := serviceFactory.MustAspectService()

	aspect = *model.NewAspect(1, 1)
	err := aspectService.Insert(&aspect)
	if err != nil {
		panic(err)
	}
}

func TestCoverServiceInsert(t *testing.T) {
	setupCoverServiceTest()
	coverService := serviceFactory.MustCoverService()
	defer coverService.Close()

	c1 := model.Cover{
		AspectId: aspect.Id,
		Width:    600,
		Height:   400,
	}

	err := coverService.Insert(&c1)
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
		c1.AspectId != c2.AspectId ||
		c1.Width != c2.Width ||
		c1.Height != c2.Height {
		t.Fatalf("Inserted cover (%+v) does not match: %+v\n", c2, c1)
	}
}

func TestCoverServiceUpdate(t *testing.T) {
	setupCoverServiceTest()
	coverService := serviceFactory.MustCoverService()
	defer coverService.Close()

	c1 := model.Cover{
		AspectId: aspect.Id,
		Width:    600,
		Height:   400,
	}

	err := coverService.Insert(&c1)
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
	setupCoverServiceTest()
	coverService := serviceFactory.MustCoverService()
	defer coverService.Close()

	c1 := model.Cover{
		AspectId: aspect.Id,
		Width:    600,
		Height:   400,
	}

	err := coverService.Insert(&c1)
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
	setupCoverServiceTest()
	coverService := serviceFactory.MustCoverService()
	defer coverService.Close()

	c1 := model.Cover{
		AspectId: aspect.Id,
		Width:    600,
		Height:   400,
	}

	err := coverService.Insert(&c1)
	if err != nil {
		t.Fatalf("Error inserting cover: %s\n", err.Error())
	}

	c2, err := coverService.GetOneBy("id", c1.Id)
	if err != nil {
		t.Fatalf("Error getting inserted cover: %s\n", err.Error())
	} else if c2 == nil {
		t.Fatalf("Cover not inserted\n")
	}

	if c1.Id != c2.Id ||
		c1.AspectId != c2.AspectId ||
		c1.Width != c2.Width ||
		c1.Height != c2.Height {
		t.Fatalf("Inserted cover (%+v) does not match: %+v\n", c2, c1)
	}
}

func TestCoverServiceFindAll(t *testing.T) {
	setupCoverServiceTest()
	coverService := serviceFactory.MustCoverService()
	defer coverService.Close()

	covers := []model.Cover{
		model.Cover{AspectId: aspect.Id, Width: 600, Height: 400},
		model.Cover{AspectId: aspect.Id, Width: 600, Height: 400},
		model.Cover{AspectId: aspect.Id, Width: 600, Height: 400},
	}

	for _, cover := range covers {
		err := coverService.Insert(&cover)
		if err != nil {
			t.Fatalf("Error inserting cover: %s\n", err.Error())
		}
	}

	covers2, err := coverService.FindAll("covers.id asc")
	if err != nil {
		t.Fatalf("Error finding covers: %s\n", err.Error())
	}

	if len(covers2) != 3 {
		t.Fatalf("Wanted 3 covers, got %d\n", len(covers2))
	}
}
