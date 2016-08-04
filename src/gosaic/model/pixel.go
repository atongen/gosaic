package model

import "encoding/json"

type Pixel interface {
	GetData() []byte
	GetPixels() []*Lab
	SetData([]byte)
	SetPixels([]*Lab)
}

// PixelEncode encodes slice of Pixels to
// json-encoded []byte and stores in Data.
func PixelEncode(p Pixel) error {
	b, err := json.Marshal(p.GetPixels())
	if err != nil {
		return err
	}
	p.SetData(b)
	return nil
}

// PixelDecode decodes []byte of Data to
// slice of *Lab and stores in Pixels.
func PixelDecode(p Pixel) error {
	var pixels []*Lab
	err := json.Unmarshal(p.GetData(), &pixels)
	if err != nil {
		return err
	}
	p.SetPixels(pixels)
	return nil
}
