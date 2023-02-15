package internal

import "strings"

func WithPrefixIfNotPresent(s, p string) string {
	if strings.HasPrefix(s, p) {
		return s
	}
	return p + s
}
