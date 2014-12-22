package mc

import "strings"

func ParseUUID(uuid string) string {
	return strings.Replace(uuid, "-", "", -1)
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
