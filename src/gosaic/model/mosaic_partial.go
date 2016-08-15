package model

type MosaicPartial struct {
	Id             int64 `db:"id"`
	MosaicId       int64 `db:"mosaic_id"`
	MacroPartialId int64 `db:"macro_partial_id"`
	GidxPartialId  int64 `db:"gidx_partial_id"`
}
