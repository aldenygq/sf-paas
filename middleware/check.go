package middleware

import "strings"

func CheckSlcHasStr(a string, b []string) bool {
	for _, v := range b {
		if strings.Contains(v, a) {
			return true
		}
	}

	return false
}
