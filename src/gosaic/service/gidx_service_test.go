package service

import (
	"testing"

	"gosaic/model"

	_ "github.com/mattn/go-sqlite3"
)

var (
	testGidx1 = &model.Gidx{
		Id:          int64(1),
		Path:        "/tmp/file.jpg",
		Md5sum:      "159c9c5ad02d9a15b7f41189960054cd",
		Width:       uint(120),
		Height:      uint(120),
		Orientation: 1,
	}
)

func setupGidxServiceTest() (GidxService, error) {
	dbMap, err := getTestDbMap()
	if err != nil {
		return nil, err
	}

	aspectService, err := getTestAspectService(dbMap)
	if err != nil {
		return nil, err
	}
	aspect = model.Aspect{Columns: 1, Rows: 1}
	err = aspectService.Insert(&aspect)
	if err != nil {
		return nil, err
	}

	gidxService, err := getTestGidxService(dbMap)
	if err != nil {
		return nil, err
	}
	testGidx1.AspectId = aspect.Id

	gidx := model.Gidx{
		AspectId:    testGidx1.AspectId,
		Path:        testGidx1.Path,
		Md5sum:      testGidx1.Md5sum,
		Width:       testGidx1.Width,
		Height:      testGidx1.Height,
		Orientation: testGidx1.Orientation,
	}
	err = gidxService.Insert(&gidx)
	if err != nil {
		return nil, err
	}
	testGidx1.Id = gidx.Id
	return gidxService, nil
}

func TestGidxServiceGet(t *testing.T) {
	gidxService, err := setupGidxServiceTest()
	if err != nil {
		t.Fatalf("Unable to setup database: %s\n", err.Error())
	}
	defer gidxService.DbMap().Db.Close()

	gidx, err := gidxService.Get(testGidx1.Id)
	if err != nil {
		t.Error("Error finding gidx by id", err)
	}

	if gidx.Id != testGidx1.Id ||
		gidx.AspectId != testGidx1.AspectId ||
		gidx.Md5sum != testGidx1.Md5sum ||
		gidx.Path != testGidx1.Path ||
		gidx.Width != testGidx1.Width ||
		gidx.Height != testGidx1.Height ||
		gidx.Orientation != testGidx1.Orientation {
		t.Error("Found gidx does not match data")
	}
}

func TestGidxServiceGetMissing(t *testing.T) {
	gidxService, err := setupGidxServiceTest()
	if err != nil {
		t.Fatalf("Unable to setup database: %s\n", err.Error())
	}
	defer gidxService.DbMap().Db.Close()

	gidx, err := gidxService.Get(1234)
	if err != nil {
		t.Error("Error finding gidx by id", err)
	}

	if gidx != nil {
		t.Error("Found non-existent gidx")
	}
}

func TestGidxServiceGetOneBy(t *testing.T) {
	gidxService, err := setupGidxServiceTest()
	if err != nil {
		t.Fatalf("Unable to setup database: %s\n", err.Error())
	}
	defer gidxService.DbMap().Db.Close()

	gidx, err := gidxService.GetOneBy("md5sum", testGidx1.Md5sum)
	if err != nil {
		t.Error("Error getting gidx for existance by md5sum", err)
	}

	if gidx.Id != testGidx1.Id ||
		gidx.AspectId != testGidx1.AspectId ||
		gidx.Md5sum != testGidx1.Md5sum ||
		gidx.Path != testGidx1.Path ||
		gidx.Width != testGidx1.Width ||
		gidx.Height != testGidx1.Height ||
		gidx.Orientation != testGidx1.Orientation {
		t.Error("Found gidx does not match data")
	}
}

func TestGidxServiceExistBy(t *testing.T) {
	gidxService, err := setupGidxServiceTest()
	if err != nil {
		t.Fatalf("Unable to setup database: %s\n", err.Error())
	}
	defer gidxService.DbMap().Db.Close()

	val, err := gidxService.ExistsBy("md5sum", testGidx1.Md5sum)
	if err != nil {
		t.Error("Error checking gidx for existance by md5sum", err)
	}

	if !val {
		t.Error("Found gidx does not exist")
	}
}

func TestGidxServiceUpdate(t *testing.T) {
	gidxService, err := setupGidxServiceTest()
	if err != nil {
		t.Fatalf("Unable to setup database: %s\n", err.Error())
	}
	defer gidxService.DbMap().Db.Close()

	newPath := "/home/user/tmp/other.jpg"
	updateGidx := model.Gidx{
		AspectId:    testGidx1.AspectId,
		Path:        newPath,
		Md5sum:      testGidx1.Md5sum,
		Width:       testGidx1.Width,
		Height:      testGidx1.Height,
		Orientation: testGidx1.Orientation,
	}
	updateGidx.Id = testGidx1.Id

	num, err := gidxService.Update(&updateGidx)
	if err != nil {
		t.Error("Error updating gidx", err)
	}

	if num == 0 {
		t.Error("Nothing was updated")
	}

	gidx, err := gidxService.Get(testGidx1.Id)
	if err != nil {
		t.Error("Error finding update gidx", err)
	}

	if gidx.Path != newPath {
		t.Error("Gidx was not updated")
	}
}

func TestGidxServiceDelete(t *testing.T) {
	gidxService, err := setupGidxServiceTest()
	if err != nil {
		t.Fatalf("Unable to setup database: %s\n", err.Error())
	}
	defer gidxService.DbMap().Db.Close()

	num, err := gidxService.Delete(&model.Gidx{Id: testGidx1.Id})
	if err != nil {
		t.Error("Error deleting gidx", err)
	}

	if num == 0 {
		t.Error("Nothing was deleted")
	}

	val, err := gidxService.ExistsBy("id", testGidx1.Id)
	if err != nil {
		t.Error("Error confirming gidx deleted", err)
	}

	if val {
		t.Error("Gidx was not deleted")
	}
}

func TestGidxServiceCount(t *testing.T) {
	gidxService, err := setupGidxServiceTest()
	if err != nil {
		t.Error("Unable to setup database", err)
	}
	defer gidxService.DbMap().Db.Close()

	num, err := gidxService.Count()
	if err != nil {
		t.Error("Error updating gidx", err)
	}

	if num == 0 {
		t.Error("Nothing was counted")
	}
}

func TestGidxServiceCountBy(t *testing.T) {
	gidxService, err := setupGidxServiceTest()
	if err != nil {
		t.Error("Unable to setup database", err)
	}
	defer gidxService.DbMap().Db.Close()

	num, err := gidxService.CountBy("md5sum", testGidx1.Md5sum)
	if err != nil {
		t.Error("Error updating gidx", err)
	}

	if num != 1 {
		t.Error("Nothing was counted")
	}
}
