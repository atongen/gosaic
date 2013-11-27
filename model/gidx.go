package model

type Gidx struct {
  Id int64
  Path string
  Md5sum string
  Width uint
  Height uint
}

func NewGidx(path string, md5sum string, width uint, height uint) *Gidx {
  return &Gidx{
    Path: path,
    Md5sum: md5sum,
    Width: width,
    Height: height}
}
