package utils

// AddStringBytes 拼接字符串, 返回 bytes from bytes.Join()
func AddStringBytes(s ...string) []byte {
	switch len(s) {
	case 0:
		return []byte{}
	case 1:
		return []byte(s[0])
	}

	n := 0
	for _, v := range s {
		n += len(v)
	}

	b := make([]byte, n)
	bp := copy(b, s[0])
	for _, v := range s[1:] {
		bp += copy(b[bp:], v)
	}

	return b
}

// AddString 拼接字符串
func AddString(s ...string) string {
	switch len(s) {
	case 0:
		return ""
	case 1:
		return s[0]
	default:
		return B2S(AddStringBytes(s...))
	}
}
