package service

import (
	"gosaic/model"
	"testing"
)

func setupMosaicServiceTest() {
	setTestDbMap()
	coverService := getTestCoverService()
	aspectService := getTestAspectService()
	macroService := getTestMacroService()

	aspect = model.Aspect{Columns: 87, Rows: 128}
	err := aspectService.Insert(&aspect)
	if err != nil {
		panic(err)
	}

	cover = model.Cover{AspectId: aspect.Id, Width: 1, Height: 1}
	cover.Name = model.CoverNameAspect(aspect.Id, 1, 1, 1)
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
}

func TestMosaicServiceInsert(t *testing.T) {
	setupMosaicServiceTest()
	mosaicService := getTestMosaicService()
	defer mosaicService.Close()

	c1 := model.Mosaic{
		Name:    "test1",
		MacroId: macro.Id,
	}

	err := mosaicService.Insert(&c1)
	if err != nil {
		t.Fatalf("Error inserting mosaic: %s\n", err.Error())
	}

	if c1.Id == int64(0) {
		t.Fatalf("Inserted mosaic id not set")
	}

	c2, err := mosaicService.Get(c1.Id)
	if err != nil {
		t.Fatalf("Error getting inserted mosaic: %s\n", err.Error())
	} else if c2 == nil {
		t.Fatalf("Mosaic not inserted\n")
	}

	if c1.Id != c2.Id ||
		c1.Name != c2.Name ||
		c1.MacroId != c2.MacroId {
		t.Fatalf("Inserted mosaic (%+v) does not match: %+v\n", c2, c1)
	}

	if c1.CreatedAt.IsZero() ||
		c2.CreatedAt.IsZero() {
		t.Fatal("Mosaic created at not set")
	} else if c1.CreatedAt.Unix() != c2.CreatedAt.Unix() {
		t.Fatal("Inserted mosaic created at does not match")
	}
}

func TestMosaicServiceGetOneBy(t *testing.T) {
	setupMosaicServiceTest()
	mosaicService := getTestMosaicService()
	defer mosaicService.Close()

	c1 := model.Mosaic{
		MacroId: macro.Id,
		Name:    "testme1",
	}

	err := mosaicService.Insert(&c1)
	if err != nil {
		t.Fatalf("Error inserting mosaic: %s\n", err.Error())
	}

	c2, err := mosaicService.GetOneBy("macro_id = ? and name = ?", macro.Id, "testme1")
	if err != nil {
		t.Fatalf("Error getting inserted mosaic: %s\n", err.Error())
	} else if c2 == nil {
		t.Fatalf("Mosaic not inserted\n")
	}

	if c1.Id != c2.Id ||
		c1.MacroId != c2.MacroId ||
		c1.Name != c2.Name ||
		c1.CreatedAt.Unix() != c2.CreatedAt.Unix() {
		t.Fatalf("Inserted mosaic (%+v) does not match: %+v\n", c2, c1)
	}
}

func TestMosaicServiceGetOneByNot(t *testing.T) {
	setupMosaicServiceTest()
	mosaicService := getTestMosaicService()
	defer mosaicService.Close()

	c, err := mosaicService.GetOneBy("macro_id = ? and name = ?", int64(123), "not a valid name")
	if err != nil {
		t.Fatalf("Error getting inserted mosaic: %s\n", err.Error())
	}

	if c != nil {
		t.Fatal("Mosaic found when should not exist")
	}
}

func TestMosaicServiceExistsBy(t *testing.T) {
	setupMosaicServiceTest()
	mosaicService := getTestMosaicService()
	defer mosaicService.Close()

	c1 := model.Mosaic{
		MacroId: macro.Id,
		Name:    "testme1",
	}

	err := mosaicService.Insert(&c1)
	if err != nil {
		t.Fatalf("Error inserting mosaic: %s\n", err.Error())
	}

	found, err := mosaicService.ExistsBy("macro_id = ? and name = ?", macro.Id, "testme1")
	if err != nil {
		t.Fatalf("Error getting inserted mosaic: %s\n", err.Error())
	} else if !found {
		t.Fatalf("Mosaic not inserted\n")
	}
}

func TestMosaicServiceExistsByNot(t *testing.T) {
	setupMosaicServiceTest()
	mosaicService := getTestMosaicService()
	defer mosaicService.Close()

	found, err := mosaicService.ExistsBy("macro_id = ? and name = ?", int64(123), "not a valid name")
	if err != nil {
		t.Fatalf("Error getting inserted mosaic: %s\n", err.Error())
	} else if found {
		t.Fatal("Mosaic found when should not exist")
	}
}

func TestMosaicServiceFindAll(t *testing.T) {
	setupMosaicServiceTest()
	mosaicService := getTestMosaicService()
	defer mosaicService.Close()

	c1 := model.Mosaic{
		Name:    "test1",
		MacroId: macro.Id,
	}

	err := mosaicService.Insert(&c1)
	if err != nil {
		t.Fatalf("Error inserting mosaic: %s\n", err.Error())
	}

	if c1.Id == int64(0) {
		t.Fatalf("Inserted mosaic id not set")
	}

	mosaics, err := mosaicService.FindAll("id asc")
	if err != nil {
		t.Fatalf("Error finding all mosaics: %s\n", err.Error())
	}

	if len(mosaics) != 1 {
		t.Fatalf("Expected 1 mosaic, got %d\n", len(mosaics))
	}

	c2 := mosaics[0]

	if c1.Id != c2.Id ||
		c1.Name != c2.Name ||
		c1.MacroId != c2.MacroId {
		t.Fatalf("Inserted mosaic (%+v) does not match: %+v\n", c2, c1)
	}

	if c1.CreatedAt.IsZero() ||
		c2.CreatedAt.IsZero() {
		t.Fatal("Mosaic created at not set")
	} else if c1.CreatedAt.Unix() != c2.CreatedAt.Unix() {
		t.Fatal("Inserted mosaic created at does not match")
	}
}
