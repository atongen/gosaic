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

func (g *Macro) GetPath() string {
	return g.Path
}

func (g *Macro) GetMd5sum() string {
	return g.Md5sum
}

func (g *Macro) GetWidth() uint {
	return g.Width
}

func (g *Macro) GetHeight() uint {
	return g.Height
}

func (g *Macro) GetOrientation() int {
	return g.Orientation
}

func (g *Macro) SetPath(path string) {
	g.Path = path
}

func (g *Macro) SetMd5sum(md5sum string) {
	g.Md5sum = md5sum
}

func (g *Macro) SetWidth(width uint) {
	g.Width = width
}

func (g *Macro) SetHeight(height uint) {
	g.Height = height
}

func (g *Macro) SetOrientation(orientation int) {
	g.Orientation = orientation
}
