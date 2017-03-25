package database

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestMigrate(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Could not get test db: %s\n", err.Error())
	}
	defer db.Close()

	version, err := Migrate(db)
	if err != nil {
		t.Fatalf("Failed to migrate db: %s\n", err.Error())
	}

	if version != len(migrations) {
		t.Fatalf("Failed to complete all migrations (%d)\n", version)
	}
}

func TestGetVersion(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Could not get test db: %s\n", err.Error())
	}
	defer db.Close()

	version, err := GetVersion(db)
	if err != nil {
		t.Fatalf("Error while getting version: %s\n", err.Error())
	}

	if version != 0 {
		t.Fatalf("Incorrect version number returned: %d\n", version)
	}
}
