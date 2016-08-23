package model

type Image interface {
	GetPath() string
	GetMd5sum() string
	GetWidth() int
	GetHeight() int
	GetOrientation() int
	SetPath(string)
	SetMd5sum(string)
	SetWidth(int)
	SetHeight(int)
	SetOrientation(int)
}
