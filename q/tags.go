package q

type tagsQuery []string

// Tags returns a tags Query
func Tags(tag string, extra ...string) Query {
	tags := make(tagsQuery, 1+len(extra))
	i := 0
	t := toLowerTrim(tag)
	if t != "" {
		tags[i] = t
		i++
	}
	for _, s := range extra {
		t := toLowerTrim(s)
		if t != "" {
			tags[i] = t
			i++
		}
	}
	return tags[:i]
}

// QueryType implements the Query interface
func (q tagsQuery) QueryType() string {
	return "tags"
}
