package model

type Gidx struct {
	Id          int64  `db:"id"`
	AspectId    int64  `db:"aspect_id"`
	Path        string `db:"path"`
	Md5sum      string `db:"md5sum"`
	Width       uint   `db:"width"`
	Height      uint   `db:"height"`
	Orientation int    `db:"orientation"`
}

func (g *Gidx) GetPath() string {
	return g.Path
}

func (g *Gidx) GetMd5sum() string {
	return g.Md5sum
}

func (g *Gidx) GetWidth() uint {
	return g.Width
}

func (g *Gidx) GetHeight() uint {
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

func (g *Gidx) SetWidth(width uint) {
	g.Width = width
}

func (g *Gidx) SetHeight(height uint) {
	g.Height = height
}

func (g *Gidx) SetOrientation(orientation int) {
	g.Orientation = orientation
}
