package botutil

import "strings"

func HasLinePrefix(s, prefix string) bool {
	pos := strings.Index(s, prefix)
	if pos > -1 && (pos == 0 || s[pos-1] == '\n') {
		return true
	}
	return false
}

func HasAnyLinePrefix(s string, prefixes []string) bool {
	for _, prefix := range prefixes {
		if HasLinePrefix(s, prefix) {
			return true
		}
	}
	return false
}
