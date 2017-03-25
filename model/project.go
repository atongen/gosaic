package model

import "time"

type Project struct {
	Id         int64     `db:"id"`
	Name       string    `db:"name"`
	Path       string    `db:"path"`
	CoverPath  string    `db:"cover_path"`
	MacroPath  string    `db:"macro_path"`
	MosaicPath string    `db:"mosaic_path"`
	CoverId    int64     `db:"cover_id"`
	MacroId    int64     `db:"macro_id"`
	MosaicId   int64     `db:"mosaic_id"`
	IsComplete bool      `db:"is_complete"`
	CreatedAt  time.Time `db:"created_at"`
}
