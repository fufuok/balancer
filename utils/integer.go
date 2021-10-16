package utils

// SearchInts similar to sort.SearchInts(), but faster.
func SearchInts(a []int, x int) (i int) {
	j := len(a)
	for i < j {
		h := int(uint(i+j) >> 1)
		if a[h] < x {
			i = h + 1
		} else {
			j = h
		}
	}
	return i
}
