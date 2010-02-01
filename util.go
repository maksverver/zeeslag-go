package util

// min returns the minimum value of its two arguments
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// max returns the maximum value of its two arguments
func Max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

// ifc returns the a or b argument, depending on whether c is true or false
func Ifc(c bool, a, b int) int {
	if c {
		return a
	}
	return b
}
