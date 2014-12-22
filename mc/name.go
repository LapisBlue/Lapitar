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

func equalLower(a, b string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		c1, c2 := a[i], b[i]
		if c1 != c2 && toLower(c1) != toLower(c2) {
			return false
		}
	}

	return true
}

func toLower(c uint8) uint8 {
	if c >= 'A' && c <= 'Z' {
		c += 'a' - 'A'
	}

	return c
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
