package service

import (
	"database/sql"
	"testing"

	"github.com/atongen/gosaic/database"
	"github.com/atongen/gosaic/model"
	_ "github.com/mattn/go-sqlite3"
)

var (
	id     int64
	path   = "/tmp/file.jpg"
	md5sum = "159c9c5ad02d9a15b7f41189960054cd"
	width  = uint(120)
	height = uint(120)
)

func setupGidxServiceTest() (*GidxService, *sql.DB, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, nil, err
	}
	_, err = database.Migrate(db)
	if err != nil {
		return nil, nil, err
	}
	gidxService := NewGidxService(db)
	gidx := model.NewGidx(path, md5sum, width, height)
	err = gidxService.Create(gidx)
	if err != nil {
		return nil, nil, err
	}
	id = gidx.Id
	return gidxService, db, nil
}

func TestGidxServiceFindById(t *testing.T) {
	gidxService, db, err := setupGidxServiceTest()
	if err != nil {
		t.Error("Unable to setup database.", err)
	}
	defer db.Close()

	gidx, err := gidxService.FindById(id)
	if err != nil {
		t.Error("Error finding gidx by id", err)
	}

	if gidx.Id != id ||
		gidx.Md5sum != md5sum ||
		gidx.Path != path ||
		gidx.Width != width ||
		gidx.Height != height {
		t.Error("Found gidx does not match data")
	}
}

func TestGidxServiceExistId(t *testing.T) {
	gidxService, db, err := setupGidxServiceTest()
	if err != nil {
		t.Error("Unable to setup database.", err)
	}
	defer db.Close()

	val, err := gidxService.ExistsById(id)
	if err != nil {
		t.Error("Error checking gidx for existance by id", err)
	}

	if !val {
		t.Error("Found gidx does not exist")
	}
}

func TestGidxServiceFindByMd5sum(t *testing.T) {
	gidxService, db, err := setupGidxServiceTest()
	if err != nil {
		t.Error("Unable to setup database.", err)
	}
	defer db.Close()

	gidx, err := gidxService.FindByMd5sum(md5sum)
	if err != nil {
		t.Error("Error finding gidx by md5sum", err)
	}

	if gidx.Id != id ||
		gidx.Md5sum != md5sum ||
		gidx.Path != path ||
		gidx.Width != width ||
		gidx.Height != height {
		t.Error("Found gidx does not match data")
	}
}

func TestGidxServiceExistByMd5sum(t *testing.T) {
	gidxService, db, err := setupGidxServiceTest()
	if err != nil {
		t.Error("Unable to setup database.", err)
	}
	defer db.Close()

	val, err := gidxService.ExistsByMd5sum(md5sum)
	if err != nil {
		t.Error("Error checking gidx for existance by md5sum", err)
	}

	if !val {
		t.Error("Found gidx does not exist")
	}
}

func TestGidxServiceUpdate(t *testing.T) {
	gidxService, db, err := setupGidxServiceTest()
	if err != nil {
		t.Error("Unable to setup database", err)
	}
	defer db.Close()

	newPath := "/home/user/tmp/other.jpg"
	updateGidx := model.NewGidx(newPath, md5sum, width, height)
	updateGidx.Id = id

	err = gidxService.Update(updateGidx)
	if err != nil {
		t.Error("Error updating gidx", err)
	}

	gidx, err := gidxService.FindById(id)
	if err != nil {
		t.Error("Error finding update gidx", err)
	}

	if gidx.Path != newPath {
		t.Error("Gidx was not updated")

	}
}

func TestGidxServiceDelete(t *testing.T) {
	gidxService, db, err := setupGidxServiceTest()
	if err != nil {
		t.Error("Unable to setup database", err)
	}
	defer db.Close()

	err = gidxService.Delete(id)
	if err != nil {
		t.Error("Error deleting gidx", err)
	}

	val, err := gidxService.ExistsById(id)
	if err != nil {
		t.Error("Error confirming gidx deleted", err)
	}

	if val {
		t.Error("Gidx was not deleted")

	}
}
