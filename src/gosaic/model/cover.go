package model

type Cover struct {
	Id     int64  `db:"id"`
	Type   string `db:"type"`
	Name   string `db:"name"`
	Width  uint   `db:"width"`
	Height uint   `db:"height"`
}
