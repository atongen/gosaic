package service

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/coopernurse/gorp"

	"github.com/atongen/gosaic/database"
	"github.com/atongen/gosaic/model"
	_ "github.com/mattn/go-sqlite3"
)

func setupGenericServiceTest() (*GenericServiceImpl, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}
	_, err = database.Migrate(db)
	if err != nil {
		return nil, err
	}
	dbMap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	dbMap.TraceOn("[DB]", log.New(os.Stdout, "test:", log.Ldate|log.Ltime))
	genericService := NewGenericService(dbMap, model.Gidx{}, "gidx", "Id")
	genericService.Register()

	gidx := model.NewGidx(path, md5sum, width, height, orientation)
	err = genericService.Insert(gidx)
	if err != nil {
		return nil, err
	}
	id = gidx.Id
	return genericService, nil
}

func TestGenericServiceGet(t *testing.T) {
	gidxService, err := setupGenericServiceTest()
	if err != nil {
		t.Error("Unable to setup database.", err)
	}
	defer gidxService.DbMap().Db.Close()

	obj, err := gidxService.Get(id)
	if err != nil {
		t.Error("Error finding gidx by id", err)
	}
	gidx := obj.(*model.Gidx)

	if gidx.Id != id ||
		gidx.Md5sum != md5sum ||
		gidx.Path != path ||
		gidx.Width != width ||
		gidx.Height != height ||
		gidx.Orientation != orientation {
		t.Error("Found gidx does not match data")
	}
}

func TestGenericServiceGetMissing(t *testing.T) {
	gidxService, err := setupGenericServiceTest()
	if err != nil {
		t.Error("Unable to setup database.", err)
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

func TestGenericServiceExistBy(t *testing.T) {
	gidxService, err := setupGenericServiceTest()
	if err != nil {
		t.Error("Unable to setup database.", err)
	}
	defer gidxService.DbMap().Db.Close()

	val, err := gidxService.ExistsBy("md5sum", md5sum)
	if err != nil {
		t.Error("Error checking gidx for existance by md5sum", err)
	}

	if !val {
		t.Error("Found gidx does not exist")
	}
}

func TestGenericServiceUpdate(t *testing.T) {
	gidxService, err := setupGenericServiceTest()
	if err != nil {
		t.Error("Unable to setup database", err)
	}
	defer gidxService.DbMap().Db.Close()

	newPath := "/home/user/tmp/other.jpg"
	updateGidx := model.NewGidx(newPath, md5sum, width, height, orientation)
	updateGidx.Id = id

	num, err := gidxService.Update(updateGidx)
	if err != nil {
		t.Error("Error updating gidx", err)
	}

	if num == 0 {
		t.Error("Nothing was updated")
	}

	obj, err := gidxService.Get(id)
	if err != nil {
		t.Error("Error finding update gidx", err)
	}
	gidx := obj.(*model.Gidx)

	if gidx.Path != newPath {
		t.Error("Gidx was not updated")
	}
}

func TestGenericServiceDelete(t *testing.T) {
	gidxService, err := setupGenericServiceTest()
	if err != nil {
		t.Error("Unable to setup database", err)
	}
	defer gidxService.DbMap().Db.Close()

	num, err := gidxService.Delete(&model.Gidx{Id: id})
	if err != nil {
		t.Error("Error deleting gidx", err)
	}

	if num == 0 {
		t.Error("Nothing was deleted")
	}

	val, err := gidxService.ExistsBy("id", id)
	if err != nil {
		t.Error("Error confirming gidx deleted", err)
	}

	if val {
		t.Error("Gidx was not deleted")
	}
}

func TestGenericServiceCount(t *testing.T) {
	gidxService, err := setupGenericServiceTest()
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

func TestGenericServiceCountBy(t *testing.T) {
	gidxService, err := setupGenericServiceTest()
	if err != nil {
		t.Error("Unable to setup database", err)
	}
	defer gidxService.DbMap().Db.Close()

	num, err := gidxService.CountBy("md5sum", md5sum)
	if err != nil {
		t.Error("Error updating gidx", err)
	}

	if num != 1 {
		t.Error("Nothing was counted")
	}
}
