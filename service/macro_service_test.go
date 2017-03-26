package service

import (
	"testing"

	"github.com/atongen/gosaic/model"

	_ "github.com/mattn/go-sqlite3"
)

func setupMacroServiceTest() {
	setTestServiceFactory()
	coverService := serviceFactory.MustCoverService()
	aspectService := serviceFactory.MustAspectService()

	aspect = model.Aspect{Columns: 1, Rows: 1}
	err := aspectService.Insert(&aspect)
	if err != nil {
		panic(err)
	}

	cover = model.Cover{AspectId: aspect.Id, Width: 1, Height: 1}
	err = coverService.Insert(&cover)
	if err != nil {
		panic(err)
	}
}

func TestMacroServiceInsert(t *testing.T) {
	setupMacroServiceTest()
	macroService := serviceFactory.MustMacroService()
	defer macroService.Close()

	c1 := model.Macro{
		AspectId:    aspect.Id,
		CoverId:     cover.Id,
		Path:        "/path/to/my/macro_image.jpg",
		Md5sum:      "68b329da9893e34099c7d8ad5cb9c940",
		Width:       1,
		Height:      1,
		Orientation: 1,
	}

	err := macroService.Insert(&c1)
	if err != nil {
		t.Fatalf("Error inserting macro: %s\n", err.Error())
	}

	if c1.Id == int64(0) {
		t.Fatalf("Inserted macro id not set")
	}

	c2, err := macroService.Get(c1.Id)
	if err != nil {
		t.Fatalf("Error getting inserted macro: %s\n", err.Error())
	} else if c2 == nil {
		t.Fatalf("Macro not inserted\n")
	}

	if c1.Id != c2.Id ||
		c1.AspectId != c2.AspectId ||
		c1.CoverId != c2.CoverId ||
		c1.Path != c2.Path ||
		c1.Md5sum != c1.Md5sum ||
		c1.Width != c2.Width ||
		c1.Height != c2.Height ||
		c1.Orientation != c2.Orientation {
		t.Fatalf("Inserted macro (%+v) does not match: %+v\n", c2, c1)
	}
}

func TestMacroServiceUpdate(t *testing.T) {
	setupMacroServiceTest()
	macroService := serviceFactory.MustMacroService()
	defer macroService.Close()

	c1 := model.Macro{
		AspectId:    aspect.Id,
		CoverId:     cover.Id,
		Path:        "/path/to/my/macro_image.jpg",
		Md5sum:      "68b329da9893e34099c7d8ad5cb9c940",
		Width:       1,
		Height:      1,
		Orientation: 1,
	}

	err := macroService.Insert(&c1)
	if err != nil {
		t.Fatalf("Error inserting macro: %s\n", err.Error())
	}

	c1.Width = 800
	err = macroService.Update(&c1)
	if err != nil {
		t.Fatalf("Error updating macro: %s\n", err.Error())
	}

	c2, err := macroService.Get(c1.Id)
	if err != nil {
		t.Fatalf("Error getting inserted macro: %s\n", err.Error())
	}

	if c2.Width != 800 {
		t.Fatalf("Error updating macro, expected width 800, got width %s\n", c2.Width)
	}
}

func TestMacroServiceDelete(t *testing.T) {
	setupMacroServiceTest()
	macroService := serviceFactory.MustMacroService()
	defer macroService.Close()

	c1 := model.Macro{
		AspectId:    aspect.Id,
		CoverId:     cover.Id,
		Path:        "/path/to/my/macro_image.jpg",
		Md5sum:      "68b329da9893e34099c7d8ad5cb9c940",
		Width:       1,
		Height:      1,
		Orientation: 1,
	}

	err := macroService.Insert(&c1)
	if err != nil {
		t.Fatalf("Error inserting macro: %s\n", err.Error())
	}

	err = macroService.Delete(&c1)
	if err != nil {
		t.Fatalf("Error deleting macro: %s\n", err.Error())
	}

	c2, err := macroService.Get(c1.Id)
	if err != nil {
		t.Fatalf("Error getting macro: %s\n", err.Error())
	} else if c2 != nil {
		t.Fatalf("Macro not deleted")
	}
}

func TestMacroServiceGetOneBy(t *testing.T) {
	setupMacroServiceTest()
	macroService := serviceFactory.MustMacroService()
	defer macroService.Close()

	c1 := model.Macro{
		AspectId:    aspect.Id,
		CoverId:     cover.Id,
		Path:        "/path/to/my/macro_image.jpg",
		Md5sum:      "68b329da9893e34099c7d8ad5cb9c940",
		Width:       1,
		Height:      1,
		Orientation: 1,
	}

	err := macroService.Insert(&c1)
	if err != nil {
		t.Fatalf("Error inserting macro: %s\n", err.Error())
	}

	c2, err := macroService.GetOneBy("md5sum = ?", "68b329da9893e34099c7d8ad5cb9c940")
	if err != nil {
		t.Fatalf("Error getting inserted macro: %s\n", err.Error())
	} else if c2 == nil {
		t.Fatalf("Macro not inserted\n")
	}

	if c1.Id != c2.Id ||
		c1.AspectId != c2.AspectId ||
		c1.CoverId != c2.CoverId ||
		c1.Path != c2.Path ||
		c1.Md5sum != c1.Md5sum ||
		c1.Width != c2.Width ||
		c1.Height != c2.Height ||
		c1.Orientation != c2.Orientation {
		t.Fatalf("Inserted macro (%+v) does not match: %+v\n", c2, c1)
	}
}

func TestMacroServiceGetOneByNone(t *testing.T) {
	setupMacroServiceTest()
	macroService := serviceFactory.MustMacroService()
	defer macroService.Close()

	c, err := macroService.GetOneBy("cover_id = ? AND md5sum = ?", int64(123), "not a valid md5")
	if err != nil {
		t.Fatalf("Error getting inserted macro: %s\n", err.Error())
	}

	if c != nil {
		t.Fatal("Macro found when should not exist")
	}
}

func TestMacroServiceExistsBy(t *testing.T) {
	setupMacroServiceTest()
	macroService := serviceFactory.MustMacroService()
	defer macroService.Close()

	c1 := model.Macro{
		AspectId:    aspect.Id,
		CoverId:     cover.Id,
		Path:        "/path/to/my/macro_image.jpg",
		Md5sum:      "68b329da9893e34099c7d8ad5cb9c940",
		Width:       1,
		Height:      1,
		Orientation: 1,
	}

	err := macroService.Insert(&c1)
	if err != nil {
		t.Fatalf("Error inserting macro: %s\n", err.Error())
	}

	found, err := macroService.ExistsBy("cover_id = ? AND md5sum = ?", cover.Id, "68b329da9893e34099c7d8ad5cb9c940")
	if err != nil {
		t.Fatalf("Error checking inserted macro: %s\n", err.Error())
	} else if !found {
		t.Fatalf("Macro not found\n")
	}
}

func TestMacroServiceFindAll(t *testing.T) {
	setupMacroServiceTest()
	macroService := serviceFactory.MustMacroService()
	defer macroService.Close()

	macros := []model.Macro{
		model.Macro{
			AspectId:    aspect.Id,
			CoverId:     cover.Id,
			Path:        "/path/to/my/macro_image.jpg",
			Md5sum:      "68b329da9893e34099c7d8ad5cb9c940",
			Width:       1,
			Height:      1,
			Orientation: 1,
		},
		model.Macro{
			AspectId:    aspect.Id,
			CoverId:     cover.Id,
			Path:        "/path/to/my/macro_image.jpg",
			Md5sum:      "68b329da9893e34099c7d8ad5cb9c941",
			Width:       1,
			Height:      1,
			Orientation: 1,
		},
		model.Macro{
			AspectId:    aspect.Id,
			CoverId:     cover.Id,
			Path:        "/path/to/my/macro_image.jpg",
			Md5sum:      "68b329da9893e34099c7d8ad5cb9c942",
			Width:       1,
			Height:      1,
			Orientation: 1,
		},
	}

	for _, macro := range macros {
		err := macroService.Insert(&macro)
		if err != nil {
			t.Fatalf("Error inserting macro: %s\n", err.Error())
		}
	}

	macros2, err := macroService.FindAll("macros.id ASC")
	if err != nil {
		t.Fatalf("Error finding macros: %s\n", err.Error())
	}

	if len(macros2) != 3 {
		t.Fatalf("Wanted 3 macros, got %d\n", len(macros2))
	}
}
