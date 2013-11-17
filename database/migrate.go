package database

import (
  "fmt"
  "database/sql"
  _ "github.com/mattn/go-sqlite3"
)

type migrationFunc func(db *sql.DB) error
var migrations []migrationFunc

func init() {
  migrations = make([]migrationFunc, 2)
  migrations[0] = createVersionTable
  migrations[1] = createGidxTable
}

func Migrate(dbPath string) error {
  // get the db
  db, err := sql.Open("sqlite3", dbPath)
  if err != nil {
    return err
  }
  defer db.Close()

  // get the current version of the db
  version, err := getVersion(db)
  if err != nil {
    return err
  }

  for idx, migFun := range migrations {
    migVer := idx + 1
    if version < migVer {
      err = migFun(db)
      if err != nil {
        return err
      }
      err = setVersion(db, migVer)
      if err != nil {
        return err
      }
      fmt.Printf("Migrated database to version %d\n", migVer)
    }
  }

  return nil
}

func getVersion(db *sql.DB) (int, error) {
  var version int
  sql := `
    select version
    from versions
    order by version desc
    limit 1
  `
  rows, err := db.Query(sql)
  if err != nil {
    // db has not been created yet
    return 0, nil
  }
  defer rows.Close()

  for rows.Next() {
          rows.Scan(&version)
  }
  rows.Close()

  return version, nil
}

func setVersion(db *sql.DB, version int) error {
  tx, err := db.Begin()
  if err != nil {
    return err
  }

  sql := "insert into versions(version) values(?)"
  stmt, err := tx.Prepare(sql)
  if err != nil {
    return err
  }
  defer stmt.Close()

  _, err = stmt.Exec(version)
  if err != nil {
    return err
  }
  tx.Commit()
  return nil
}

func createVersionTable(db *sql.DB) error {
  sql := `
    create table versions (
      version integer not null primary key
    );
  `
  _, err := db.Exec(sql)

  if err != nil {
    return err
  }

  sql = "create unique index idx_versions_version on versions (version)"
  _, err = db.Exec(sql)
  return err
}

func createGidxTable(db *sql.DB) error {
  sql := `
    create table gidx (
      id integer not null primary key,
      path text not null,
      md5sum text not null,
      width integer not null,
      height integer not null
    );
  `
  _, err := db.Exec(sql)
  if err != nil {
    return err
  }

  sql = "create unique index idx_gidx_md5sum on gidx (md5sum);"
  _, err = db.Exec(sql)
  return err
}
