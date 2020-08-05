package q

// Resolver to save query resolves
type Resolver struct {
	IDs     []string
	Filters []func(Post) bool
	Tags    []string
	Keys    []string
	Owner   string
	Anchor  string
	Offset  string
	Limit   uint32
	Order   uint8
}

// Filter checks all the filters
func (res Resolver) Filter(p Post) bool {
	for _, f := range res.Filters {
		if !f(p) {
			return false
		}
	}
	return true
}
