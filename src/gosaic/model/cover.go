package model

type Cover struct {
	Id       int64  `db:"id"`
	AspectId int64  `db:"aspect_id"`
	Type     string `db:"type"`
	Width    uint   `db:"width"`
	Height   uint   `db:"height"`
}
