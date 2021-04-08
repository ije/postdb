package postdb

// prefix creates a prefix
func toPrefix(a string) []byte {
	buf := make([]byte, 1+len(a))
	copy(buf, a)
	return buf
}

// join joins a and b by char c
func join(a []byte, b []byte, c byte) []byte {
	l := len(a)
	buf := make([]byte, l+1+len(b))
	copy(buf, a)
	buf[l] = c
	copy(buf[l+1:], b)
	return buf
}

// containsSlice reports whether b is within a.
func containsSlice(a []string, b []string) bool {
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
