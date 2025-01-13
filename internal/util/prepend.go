package util

// from https://stackoverflow.com/a/41918820
func Prepend[T any](dest []T, value T) []T {
	if cap(dest) > len(dest) {
		dest = dest[:len(dest)+1]
		copy(dest[1:], dest)
		dest[0] = value
		return dest
	}

	// No room, new slice need to be allocated:
	// Use some extra space for future:
	res := make([]T, len(dest)+1, len(dest)+5)
	res[0] = value
	copy(res[1:], dest)
	return res
}
