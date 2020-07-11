package q

import (
	"strings"
)

func toLowerTrim(s string) string {
	return strings.ToLower(strings.TrimSpace(string(s)))
}
