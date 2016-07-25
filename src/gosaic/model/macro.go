package model

type Macro struct {
	Id          int64  `db:"id"`
	AspectId    int64  `db:"aspect_id"`
	CoverId     int64  `db:"cover_id"`
	Path        string `db:"path"`
	Md5sum      string `db:"md5sum"`
	Width       uint   `db:"width"`
	Height      uint   `db:"height"`
	Orientation int    `db:"orientation"`
}
