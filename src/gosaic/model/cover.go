package model

type Cover struct {
	Id     int64  `db:"id"`
	Name   string `db:"name"`
	Width  uint   `db:"width"`
	Height uint   `db:"height"`
}
