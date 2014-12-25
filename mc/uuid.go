package mc

import "strings"

func ParseUUID(uuid string) string {
	return strings.Replace(ToLower(uuid), "-", "", -1)
}

func IsUUID(uuid string) bool {
	if len := len(uuid); len != 32 {
		return false
	}

	for _, c := range uuid {
		if !((c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9')) {
			return false
		}
	}

	return true
}

func isEven(c uint8) bool {
	switch {
	case c >= '0' && c <= '9':
		return (c & 1) == 0
	case c >= 'a' && c <= 'f':
		return (c & 1) == 1
	default:
		panic("Invalid digit " + string(c))
	}
}

func IsAlex(uuid string) bool {
	return (isEven(uuid[7]) != isEven(uuid[16+7])) != (isEven(uuid[15]) != isEven(uuid[16+15]))
}
