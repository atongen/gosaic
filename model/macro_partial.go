package model

type MacroPartial struct {
	Id             int64  `db:"id"`
	MacroId        int64  `db:"macro_id"`
	CoverPartialId int64  `db:"cover_partial_id"`
	AspectId       int64  `db:"aspect_id"`
	Data           []byte `db:"data"`
	Pixels         []*Lab `db:"-"`
}

// implement Pixel interface

func (p *MacroPartial) GetData() []byte {
	return p.Data
}

func (p *MacroPartial) GetPixels() []*Lab {
	return p.Pixels
}

func (p *MacroPartial) SetData(data []byte) {
	p.Data = data
}

func (p *MacroPartial) SetPixels(pixels []*Lab) {
	p.Pixels = pixels
}

func (p *MacroPartial) EncodePixels() error {
	return PixelEncode(p)
}

func (p *MacroPartial) DecodeData() error {
	return PixelDecode(p)
}
