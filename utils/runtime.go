package utils

import (
	_ "unsafe"
)

// FastRandn similar to fastrand() % n, but faster.
// See https://lemire.me/blog/2016/06/27/a-fast-alternative-to-the-modulo-reduction/
//go:linkname FastRandn runtime.fastrandn
func FastRandn(n uint32) uint32
