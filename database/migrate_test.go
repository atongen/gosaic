package database

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestMigrate(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Error("Could not get test db.", err)
	}
	defer db.Close()
	version, err := Migrate(db)
	if err != nil {
		t.Error("Failed to migrate db.", err)
	}
	if version != 3 {
		t.Error("Failed to complete all migrations.", version)
	}
}

func TestGetVersion(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Error("Could not get test db.", err)
	}
	defer db.Close()
	version, err := GetVersion(db)
	if err != nil {
		t.Error("Errro while getting version", err)
	}
	if version != 0 {
		t.Error("Incorrect version number returned.", version)
	}
}
