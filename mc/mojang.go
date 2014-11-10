package mc

func IsName(name string) bool {
	if len := len(name); len == 0 || len > 16 {
		return false
	}

	for _, c := range name {
		if !((c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '_') {
			return false
		}
	}

	return true
}

func ToLower(name string) string {
	buf := make([]rune, len(name))
	for i, c := range name {
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}

		buf[i] = c
	}

	return string(buf)
}
