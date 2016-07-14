package service

import (
	"database/sql"
	"gosaic/database"

	"gopkg.in/gorp.v1"
)

func getTestDbMap() (*gorp.DbMap, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return nil, err
	}

	_, err = database.Migrate(db)
	if err != nil {
		return nil, err
	}
	dbMap := gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}

	return &dbMap, nil
}
