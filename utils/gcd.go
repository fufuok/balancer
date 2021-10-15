package utils

// GCD Greatest Common Divisor
func GCD(x, y int) int {
	for y != 0 {
		x, y = y, x%y
	}
	return x
}
