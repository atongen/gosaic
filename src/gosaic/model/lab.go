package model

import (
	"image/color"
	"math"

	"github.com/lucasb-eyer/go-colorful"
)

type Lab struct {
	L     float64 `json:"l"`
	A     float64 `json:"a"`
	B     float64 `json:"b"`
	Alpha float64 `json:"alpha"`
}

func (lab1 *Lab) Dist(lab2 *Lab) float64 {
	return math.Sqrt(sq(lab1.L-lab2.L) + sq(lab1.A-lab2.A) + sq(lab1.B-lab2.B))
}

func sq(v float64) float64 {
	return v * v
}

func RgbaToLab(color color.Color) *Lab {
	r, g, b, alpha := color.RGBA()
	myColor := colorful.Color{R: float64(r) / 65535.0, G: float64(g) / 65535.0, B: float64(b) / 65535.0}
	l, a, bb := myColor.Lab()
	return &Lab{l, a, bb, float64(alpha)}
}
