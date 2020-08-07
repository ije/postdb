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
