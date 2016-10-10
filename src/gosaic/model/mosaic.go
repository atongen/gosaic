package model

type Mosaic struct {
	Id      int64 `db:"id"`
	MacroId int64 `db:"macro_id"`
}
