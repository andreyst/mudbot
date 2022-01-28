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

func HasLineSuffix(s, suffix string) bool {
	pos := strings.LastIndex(s, suffix)
	if pos > -1 && (pos+len(suffix) == len(s) || s[pos+len(suffix)] == '\r') {
		return true
	}
	return false
}

func HasAnyLineSuffix(s string, suffixes []string) bool {
	for _, suffix := range suffixes {
		if HasLineSuffix(s, suffix) {
			return true
		}
	}
	return false
}
