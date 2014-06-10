package model

type Gidx struct {
	Id          int64  `db:"id"`
	AspectId    int64  `db:"aspect_id"`
	Path        string `db:"path"`
	Md5sum      string `db:"md5sum"`
	Width       uint   `db:"width"`
	Height      uint   `db:"height"`
	Orientation string `db:"orientation"`
}

func NewGidx(aspect_id int64, path string, md5sum string, width uint, height uint, orientation string) *Gidx {
	return &Gidx{
		AspectId:    aspect_id,
		Path:        path,
		Md5sum:      md5sum,
		Width:       width,
		Height:      height,
		Orientation: orientation}
}

func GidxsToInterface(gidxs []*Gidx) []interface{} {
	n := len(gidxs)
	interfaces := make([]interface{}, n)
	for i := 0; i < n; i++ {
		interfaces[i] = interface{}(gidxs[i])
	}
	return interfaces
}
