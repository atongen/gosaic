package model

import "math"

type Gidx struct {
	Id          int64  `db:"id"`
	AspectId    int64  `db:"aspect_id"`
	Path        string `db:"path"`
	Md5sum      string `db:"md5sum"`
	Width       int    `db:"width"`
	Height      int    `db:"height"`
	Orientation int    `db:"orientation"`
}

func (gidx *Gidx) Within(threashold float64, aspect *Aspect) bool {
	gRatio := float64(gidx.Width) / float64(gidx.Height)
	aRatio := aspect.Ratio()

	if gRatio == aRatio {
		return true
	}

	return math.Abs(gRatio-aRatio) <= threashold
}

// implement Image interface

func (g *Gidx) GetPath() string {
	return g.Path
}

func (g *Gidx) GetMd5sum() string {
	return g.Md5sum
}

func (g *Gidx) GetWidth() int {
	return g.Width
}

func (g *Gidx) GetHeight() int {
	return g.Height
}

func (g *Gidx) GetOrientation() int {
	return g.Orientation
}

func (g *Gidx) SetPath(path string) {
	g.Path = path
}

func (g *Gidx) SetMd5sum(md5sum string) {
	g.Md5sum = md5sum
}

func (g *Gidx) SetWidth(width int) {
	g.Width = width
}

func (g *Gidx) SetHeight(height int) {
	g.Height = height
}

func (g *Gidx) SetOrientation(orientation int) {
	g.Orientation = orientation
}
