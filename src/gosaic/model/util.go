package model

func gcd(a, b int) int {
	c := a % b

	if c == 0 {
		return b
	}

	return gcd(b, c)
}
