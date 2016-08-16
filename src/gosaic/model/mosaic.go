package model

import "time"

type Mosaic struct {
	Id        int64     `db:"id"`
	Name      string    `db:"name"`
	MacroId   int64     `db:"macro_id"`
	CreatedAt time.Time `db:"created_at"`
}
