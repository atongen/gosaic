package model

type Gidx struct {
	Id          int64
	Path        string
	Md5sum      string
	Width       uint
	Height      uint
	Orientation string
}

func NewGidx(path string, md5sum string, width uint, height uint, orientation string) *Gidx {
	return &Gidx{
		Path:        path,
		Md5sum:      md5sum,
		Width:       width,
		Height:      height,
		Orientation: orientation}
}

func (gidx *Gidx) GetId() int64 {
	return gidx.Id
}

func GidxsToInterface(gidxs []*Gidx) []interface{} {
	n := len(gidxs)
	interfaces := make([]interface{}, n)
	for i := 0; i < n; i++ {
		interfaces[i] = interface{}(gidxs[i])
	}
	return interfaces
}
