package botutil

import (
	"regexp"
	"strings"
)

var ansiRe = regexp.MustCompile(`\x1B\[\d+(;\d+)?m`)

func StripAnsi(s string) (out string) {
	if !strings.Contains(s, "\x1b") {
		out = s
		return
	}

	out = ansiRe.ReplaceAllString(s, "")
	return
}
