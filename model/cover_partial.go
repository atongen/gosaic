package model

import "image"

type CoverPartial struct {
	Id       int64 `db:"id"`
	CoverId  int64 `db:"cover_id"`
	AspectId int64 `db:"aspect_id"`
	X1       int   `db:"x1"`
	Y1       int   `db:"y1"`
	X2       int   `db:"x2"`
	Y2       int   `db:"y2"`
}

func (cp *CoverPartial) Rectangle() image.Rectangle {
	return image.Rect(int(cp.X1), int(cp.Y1), int(cp.X2), int(cp.Y2))
}

func (cp *CoverPartial) Pt() image.Point {
	return image.Point{int(cp.X1), int(cp.Y1)}
}

func (cp *CoverPartial) Width() int {
	return int(cp.X2 - cp.X1)
}

func (cp *CoverPartial) Height() int {
	return int(cp.Y2 - cp.Y1)
}

func (cp *CoverPartial) Area() int {
	return cp.Width() * cp.Height()
}
