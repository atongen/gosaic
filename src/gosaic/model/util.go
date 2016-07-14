package model

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
