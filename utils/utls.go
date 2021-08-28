package utils

// prefix creates a prefix
func ToPrefix(a string) []byte {
	buf := make([]byte, 1+len(a))
	copy(buf, a)
	return buf
}

// join joins a and b by char c
func Join(a []byte, b []byte, c byte) []byte {
	l := len(a)
	buf := make([]byte, l+1+len(b))
	copy(buf, a)
	buf[l] = c
	copy(buf[l+1:], b)
	return buf
}

// contains reports whether b is within a.
func Contains(a []string, b []string) bool {
	if len(a) < len(b) {
		return false
	}
	for _, s := range a {
		for _, _s := range b {
			if s != _s {
				return false
			}
		}
	}
	return true
}

func NoRepeat(s []string) []string {
	set := map[string]struct{}{}
	a := make([]string, len(s))
	i := 0
	for _, item := range s {
		if item != "" {
			_, ok := set[item]
			if !ok {
				set[item] = struct{}{}
				a[i] = item
				i++
			}
		}
	}
	return a[:i]
}
