package q

// StringSet returns a string set
func StringSet(arr []string) []string {
	set := map[string]struct{}{}
	a := make([]string, len(arr))
	i := 0
	for _, item := range arr {
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

// Contains reports whether s is within arr.
func Contains(arr []string, s string) bool {
	for _, item := range arr {
		if item == s {
			return true
		}
	}
	return false
}

// ContainsSlice reports whether s is within arr.
func ContainsSlice(arr []string, s []string) bool {
	for _, item := range arr {
		for _, _item := range s {
			if item != _item {
				return false
			}
		}
	}
	return true
}
