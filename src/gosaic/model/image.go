package model

type Image interface {
	GetPath() string
	GetMd5sum() string
	GetWidth() uint
	GetHeight() uint
	GetOrientation() int
	SetPath(string)
	SetMd5sum(string)
	SetWidth(uint)
	SetHeight(uint)
	SetOrientation(int)
}
