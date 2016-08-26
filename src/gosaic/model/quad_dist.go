package model

type QuadDist struct {
	Id             int64   `db:"id"`
	MacroPartialId int64   `db:"macro_partial_id"`
	Dist           float64 `db:"dist"`
}
