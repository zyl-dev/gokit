package mathutils

// MaxUint64 returns the larger of a and b.
func MaxUint64(a, b uint64) uint64 {
	if a > b {
		return a
	}

	return b
}

// MinUint64 returns the smaller of a and b.
func MinUint64(a, b uint64) uint64 {
	if a < b {
		return a
	}

	return b
}
