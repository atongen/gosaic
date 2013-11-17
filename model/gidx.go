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
  return &NewGidx{
    Path: path,
    Md5sum: md5sum,
    Width: width,
    Height: height
  }
}

func FindGidxById(db *sql.DB, id uint) *Gidx {

}

func FindGidxByMd5sum(db *sql.DB, md5sum string) *Gidx {

}

func GidxExistsByMd5sum(db *sql.DB, md5sum string) bool {

}

//
// Instance methods
//

func (gidx *Gidx) GetId() uint {
  return gidx.Id
}

func (gidx *Gidx) Create(db *sql.DB) bool {

}

func (gidx *Gidx) Update(db *sql.DB) bool {

}

func (gidx *Gidx) Delete(db *sql.DB) bool {

}
