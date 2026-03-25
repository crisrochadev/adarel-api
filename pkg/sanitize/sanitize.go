package sanitize

import (
	"regexp"
	"strings"
)

var htmlTagRegex = regexp.MustCompile(`<[^>]*>`)

func Text(input string) string {
	trimmed := strings.TrimSpace(input)
	return htmlTagRegex.ReplaceAllString(trimmed, "")
}
