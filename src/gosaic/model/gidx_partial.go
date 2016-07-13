package model

import "encoding/json"

type GidxPartial struct {
	Id       int64  `db:"id"`
	GidxId   int64  `db:"gidx_id"`
	AspectId int64  `db:"aspect_id"`
	Data     []byte `db:"data"`
	Pixels   []*Lab
}

func NewGidxPartial(gidx_id, aspect_id int64, pixels []*Lab) *GidxPartial {
	return &GidxPartial{
		GidxId:   gidx_id,
		AspectId: aspect_id,
		Pixels:   pixels,
	}
}

// EncodePixels encodes slice of Pixels to
// json-encoded []byte and stores in Data.
func (p *GidxPartial) EncodePixels() error {
	b, err := json.Marshal(p.Pixels)
	if err != nil {
		return err
	}
	p.Data = b
	return nil
}

// DecodeData decodes []byte of Data to
// slice of *Lab and stores in Pixels.
func (p *GidxPartial) DecodeData() error {
	var pixels []*Lab
	err := json.Unmarshal(p.Data, pixels)
	if err != nil {
		return err
	}
	p.Pixels = pixels
	return nil
}

func GidxPartialsToInterface(gidxPartials []*GidxPartial) []interface{} {
	n := len(gidxPartials)
	interfaces := make([]interface{}, n)
	for i := 0; i < n; i++ {
		interfaces[i] = interface{}(gidxPartials[i])
	}
	return interfaces
}
