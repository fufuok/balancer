package utils

const (
	// FNVa offset basis. See https://en.wikipedia.org/wiki/Fowler–Noll–Vo_hash_function#FNV-1a_hash
	offset64 = 14695981039346656037
	prime64  = 1099511628211
)

func HashString(s ...string) uint64 {
	return Sum64(AddString(s...))
}

// Sum64 similar to fnv.New64a().Sum64(), but faster.
func Sum64(s string) uint64 {
	var h uint64 = offset64
	for i := 0; i < len(s); i++ {
		h ^= uint64(s[i])
		h *= prime64
	}
	return h
}
