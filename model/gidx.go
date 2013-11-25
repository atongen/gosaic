package model

import (
  "database/sql"
)

type Gidx struct {
  Id uint
  Path string
  Md5sum string
  Width uint
  Height uint
}

//
// Repository methods
//

func NewGidx(path string, md5sum string, width uint, height uint) *Gidx {
  return &Gidx{
    Path: path,
    Md5sum: md5sum,
    Width: width,
    Height: height
  }
}

func FindGidxById(db *sql.DB, id uint) (*Gidx, error) {
  gidx := &Gidx{ Id: id }
  rows, err := db.Query("select path, md5sum, width, height from gidx where id = ? limit 1", id)
  if err != nil {
      return nil, err
  }
  defer rows.Close()
  for rows.Next() {
      err := rows.Scan(&gidx.Path, &gidx.Md5sum, &gidx.Width, &gidx.Height)
      if err != nil {
          return nil, err
      }
  }
  err = rows.Err()
  if err != nil {
      return gidx, err
  }
  return gidx, nil
}

func FindGidxByMd5sum(db *sql.DB, md5sum string) (*Gidx, error) {
  gidx := &Gidx{ Md5sum: md5sum }
  rows, err := db.Query("select id, path, width, height from gidx where md5sum = ? limit 1", md5sum)
  if err != nil {
      return nil, err
  }
  defer rows.Close()
  for rows.Next() {
      err := rows.Scan(&gidx.Id, &gidx.Path, &gidx.Width, &gidx.Height)
      if err != nil {
          return nil, err
      }
  }
  err = rows.Err()
  if err != nil {
      return gidx, err
  }
  return gidx, nil
}

func GidxExistsByMd5sum(db *sql.DB, md5sum string) (bool, error) {
  var exists;
  rows, err := db.Query("select 1 from gidx where md5sum = ? limit 1", md5sum)
  if err != nil {
      return nil, err
  }
  defer rows.Close()
  for rows.Next() {
      err := rows.Scan(&exists)
      if err != nil {
          return nil, err
      }
  }
  err = rows.Err()
  if err != nil {
      return nil, err
  }
  if exists == 1 {
    return true, nil
  } else {
    return false, nil
  }

}

//
// Instance methods
//

func (gidx *Gidx) GetId() uint {
  return gidx.Id
}

func (gidx *Gidx) Create(db *sql.DB) error {
  stmt, err := db.Prepare("INSERT INTO gidx(path, md5sum, width, height) VALUES(?, ?, ?, ?)")
  if err != nil {
      return err
  }
  res, err := stmt.Exec(gidx.Path, gidx.Md5sum, gidx.Width, gidx.Height)
  if err != nil {
      return err
  }
  lastId, err := res.LastInsertId()
  if err != nil {
      return err
  }
  gidx.Id = lastId
  return nil
}

func (gidx *Gidx) Update(db *sql.DB) error {
  stmt, err := db.Prepare("UPDATE gidx set path = ?, md5sum = ?, width = ?, height = ? where id = ? limit 1")
  if err != nil {
      return err
  }
  res, err := stmt.Exec(gidx.Path, gidx.Md5sum, gidx.Width, gidx.Height, gidx.Id)
  if err != nil {
      return err
  }
  return nil
}

func (gidx *Gidx) Delete(db *sql.DB) bool {
  stmt, err := db.Prepare("delete from gidx where id = ? limit 1")
  if err != nil {
      return err
  }
  res, err := stmt.Exec(gidx.Id)
  if err != nil {
      return err
  }
  return nil
}
