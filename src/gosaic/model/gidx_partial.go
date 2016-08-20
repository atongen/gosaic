package model

type GidxPartial struct {
	Id       int64  `db:"id"`
	GidxId   int64  `db:"gidx_id"`
	AspectId int64  `db:"aspect_id"`
	Data     []byte `db:"data"`
	Pixels   []*Lab `db:"-"`
}

// implement Pixel interface

func (p *GidxPartial) GetData() []byte {
	return p.Data
}

func (p *GidxPartial) GetPixels() []*Lab {
	return p.Pixels
}

func (p *GidxPartial) SetData(data []byte) {
	p.Data = data
}

func (p *GidxPartial) SetPixels(pixels []*Lab) {
	p.Pixels = pixels
}

func (p *GidxPartial) EncodePixels() error {
	return PixelEncode(p)
}

func (p *GidxPartial) DecodeData() error {
	return PixelDecode(p)
}
