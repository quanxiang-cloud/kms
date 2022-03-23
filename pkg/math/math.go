package math

// MinInt min(a,b) int
func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// MaxInt max(a, b) int
func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
