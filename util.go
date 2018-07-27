package main

import (
	"regexp"
	"strings"
)

var (
	linkRegexp = regexp.MustCompile(`<((mailto:)?[@\w\:\/\.]+)\|([@\w\:\/\.]+)>`)
)

func cleanupText(text string) string {
	matches := linkRegexp.FindAllStringSubmatch(text, -1)
	for _, m := range matches {
		text = strings.Replace(text, m[0], m[3], 1)
	}
	return text
}
