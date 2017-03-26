package service

import (
	"testing"

	"github.com/atongen/gosaic/model"

	_ "github.com/mattn/go-sqlite3"
)

func setupGidxServiceTest() {
	setTestServiceFactory()
	aspectService := serviceFactory.MustAspectService()
	gidxService := serviceFactory.MustGidxService()

	aspect = model.Aspect{Columns: 1, Rows: 1}
	err := aspectService.Insert(&aspect)
	if err != nil {
		panic(err)
	}

	gidx = model.Gidx{
		AspectId:    aspect.Id,
		Path:        "/tmp/file.jpg",
		Md5sum:      "159c9c5ad02d9a15b7f41189960054cd",
		Width:       120,
		Height:      120,
		Orientation: 1,
	}

	err = gidxService.Insert(&gidx)
	if err != nil {
		panic(err)
	}
}

func TestGidxServiceGet(t *testing.T) {
	setupGidxServiceTest()
	gidxService := serviceFactory.MustGidxService()
	defer gidxService.Close()

	gidx2, err := gidxService.Get(gidx.Id)
	if err != nil {
		t.Error("Error finding gidx by id", err)
	}

	if gidx.Id != gidx2.Id ||
		gidx.AspectId != gidx2.AspectId ||
		gidx.Md5sum != gidx2.Md5sum ||
		gidx.Path != gidx2.Path ||
		gidx.Width != gidx2.Width ||
		gidx.Height != gidx2.Height ||
		gidx.Orientation != gidx2.Orientation {
		t.Error("Found gidx does not match data")
	}
}

func TestGidxServiceGetMissing(t *testing.T) {
	setupGidxServiceTest()
	gidxService := serviceFactory.MustGidxService()
	defer gidxService.Close()

	gidx2, err := gidxService.Get(1234)
	if err != nil {
		t.Error("Error finding gidx by id", err)
	}

	if gidx2 != nil {
		t.Error("Found non-existent gidx")
	}
}

func TestGidxServiceGetOneBy(t *testing.T) {
	setupGidxServiceTest()
	gidxService := serviceFactory.MustGidxService()
	defer gidxService.Close()

	gidx2, err := gidxService.GetOneBy("md5sum", gidx.Md5sum)
	if err != nil {
		t.Error("Error getting gidx for existance by md5sum", err)
	}

	if gidx.Id != gidx2.Id ||
		gidx.AspectId != gidx2.AspectId ||
		gidx.Md5sum != gidx2.Md5sum ||
		gidx.Path != gidx2.Path ||
		gidx.Width != gidx2.Width ||
		gidx.Height != gidx2.Height ||
		gidx.Orientation != gidx2.Orientation {
		t.Error("Found gidx does not match data")
	}
}

func TestGidxServiceExistBy(t *testing.T) {
	setupGidxServiceTest()
	gidxService := serviceFactory.MustGidxService()
	defer gidxService.Close()

	val, err := gidxService.ExistsBy("md5sum", gidx.Md5sum)
	if err != nil {
		t.Error("Error checking gidx for existance by md5sum", err)
	}

	if !val {
		t.Error("Found gidx does not exist")
	}
}

func TestGidxServiceUpdate(t *testing.T) {
	setupGidxServiceTest()
	gidxService := serviceFactory.MustGidxService()
	defer gidxService.Close()

	newPath := "/home/user/tmp/other.jpg"
	updateGidx := model.Gidx{
		Id:          gidx.Id,
		AspectId:    gidx.AspectId,
		Path:        newPath,
		Md5sum:      gidx.Md5sum,
		Width:       gidx.Width,
		Height:      gidx.Height,
		Orientation: gidx.Orientation,
	}

	num, err := gidxService.Update(&updateGidx)
	if err != nil {
		t.Error("Error updating gidx", err)
	}

	if num == 0 {
		t.Error("Nothing was updated")
	}

	gidx2, err := gidxService.Get(updateGidx.Id)
	if err != nil {
		t.Error("Error finding update gidx", err)
	}

	if gidx2.Path != newPath {
		t.Error("Gidx was not updated")
	}
}

func TestGidxServiceDelete(t *testing.T) {
	setupGidxServiceTest()
	gidxService := serviceFactory.MustGidxService()
	defer gidxService.Close()

	num, err := gidxService.Delete(&model.Gidx{Id: gidx.Id})
	if err != nil {
		t.Error("Error deleting gidx", err)
	}

	if num == 0 {
		t.Error("Nothing was deleted")
	}

	val, err := gidxService.ExistsBy("id", gidx.Id)
	if err != nil {
		t.Error("Error confirming gidx deleted", err)
	}

	if val {
		t.Error("Gidx was not deleted")
	}
}

func TestGidxServiceCount(t *testing.T) {
	setupGidxServiceTest()
	gidxService := serviceFactory.MustGidxService()
	defer gidxService.Close()

	num, err := gidxService.Count()
	if err != nil {
		t.Errorf("Error updating gidx", err)
	}

	if num != 1 {
		t.Errorf("Expected 1 gidx, but got %d\n", num)
	}
}

func TestGidxServiceCountBy(t *testing.T) {
	setupGidxServiceTest()
	gidxService := serviceFactory.MustGidxService()
	defer gidxService.Close()

	num, err := gidxService.CountBy("md5sum", gidx.Md5sum)
	if err != nil {
		t.Error("Error counting by gidx", err)
	}

	if num != 1 {
		t.Errorf("Expected 1 gidx, but got %d\n", num)
	}
}
