package postdb

import (
	"strings"
)

func t2oLowerTrim(s string) string {
	return strings.ToLower(strings.TrimSpace(string(s)))
}
