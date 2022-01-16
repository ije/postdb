package q

import "github.com/ije/postdb/post"

// Resolver to save query resolves
type Resolver struct {
	IDs     []string
	Filters []func(post.Post) bool
	Tags    []string
	Keys    []string
	Owner   string
	Anchor  string
	Offset  uint32
	Limit   uint32
	Order   uint8
}

// Filter checks all the filters
func (res Resolver) Filter(p post.Post) bool {
	for _, f := range res.Filters {
		if !f(p) {
			return false
		}
	}
	return true
}
