package model

import "math"

func calculateAspect(w, h int) (int, int) {
	var d int = gcd(w, h)
	return w / d, h / d
}

func gcd(a, b int) int {
	c := a % b

	if c == 0 {
		return b
	}

	return gcd(b, c)
}

func round(f float64) int {
	var r float64
	if f >= float64(0.0) {
		r = math.Floor(f + 0.5)
	} else {
		r = math.Ceil(f - 0.5)
	}
	return int(r)
}
