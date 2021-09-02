package feed

import (
	"strings"
)

func Contain(s string, filter []string) bool {
	for _, f := range filter {
		if strings.Contains(s, f) {
			return true
		}
	}
	return false
}
