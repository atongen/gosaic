package model

import(
  "fmt"
)

type Mindex struct (
  path string
  md5sum string
  width uint
  height uint
)

func (mindex Mindex) TableName() string {
  return "mindex"
}
