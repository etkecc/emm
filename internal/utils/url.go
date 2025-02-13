package utils

import (
	"regexp"
	"strings"
)

var alphanumeric = regexp.MustCompile(`[^a-z0-9_-]+`)

// MakeURLSafe makes a string URL safe
func MakeURLSafe(s string) string {
	s = strings.TrimSpace(strings.ToLower(s))
	s = alphanumeric.ReplaceAllString(s, "-")

	// Remove leading and trailing hyphens
	for strings.HasPrefix(s, "-") {
		s = strings.TrimPrefix(s, "-")
	}
	for strings.HasSuffix(s, "-") {
		s = strings.TrimSuffix(s, "-")
	}

	return s
}
