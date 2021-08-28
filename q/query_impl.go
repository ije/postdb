package q

import (
	"github.com/postui/postdb/post"
	"github.com/postui/postdb/utils"
)

type idQuery string

// Apply implements the Query interface
func (q idQuery) Apply(p *post.Post) {}

// Resolve implements the Query interface
func (q idQuery) Resolve(r *Resolver) {
	r.IDs = append(r.IDs, string(q))
}

type idsQuery []string

// Apply implements the Query interface
func (q idsQuery) Apply(p *post.Post) {}

// Resolve implements the Query interface
func (q idsQuery) Resolve(r *Resolver) {
	r.IDs = append(r.IDs, q...)
}

type aliasQuery string

// Apply implements the Query interface
func (q aliasQuery) Apply(p *post.Post) {
	p.Alias = string(q)
}

// Resolve implements the Query interface
func (q aliasQuery) Resolve(r *Resolver) {
	r.IDs = append(r.IDs, string(q))
}

type ownerQuery string

// Apply implements the Query interface
func (q ownerQuery) Apply(p *post.Post) {
	p.Owner = string(q)
}

// Resolve implements the Query interface
func (q ownerQuery) Resolve(r *Resolver) {
	r.Owner = string(q)
}

type statusQuery uint8

// Apply implements the Query interface
func (q statusQuery) Apply(p *post.Post) {
	p.Status = uint8(q)
}

// Resolve implements the Query interface
func (q statusQuery) Resolve(r *Resolver) {
	r.Filters = append(r.Filters, func(p post.Post) bool {
		return p.Status == uint8(q)
	})
}

type tagsQuery []string

// Apply implements the Query interface
func (q tagsQuery) Apply(p *post.Post) {
	p.Tags = q
}

// Resolve implements the Query interface
func (q tagsQuery) Resolve(r *Resolver) {
	tags := append(r.Tags, q...)
	r.Tags = utils.NoRepeat(tags)
}

type keysQuery []string

// Apply implements the Query interface
func (q keysQuery) Apply(p *post.Post) {}

// Resolve implements the Query interface
func (q keysQuery) Resolve(r *Resolver) {
	keys := append(r.Keys, q...)
	r.Keys = utils.NoRepeat(keys)
}

type anchorQuery string

// Apply implements the Query interface
func (q anchorQuery) Apply(p *post.Post) {}

// Resolve implements the Query interface
func (q anchorQuery) Resolve(r *Resolver) {
	r.Anchor = string(q)
}

type offsetQuery uint32

// Apply implements the Query interface
func (q offsetQuery) Apply(p *post.Post) {}

// Resolve implements the Query interface
func (q offsetQuery) Resolve(r *Resolver) {
	r.Offset = uint32(q)
}

type limitQuery uint32

// Apply implements the Query interface
func (q limitQuery) Apply(p *post.Post) {}

// Resolve implements the Query interface
func (q limitQuery) Resolve(r *Resolver) {
	r.Limit = uint32(q)
}

type rangeQuery [2]uint32

// Apply implements the Query interface
func (q rangeQuery) Apply(p *post.Post) {}

// Resolve implements the Query interface
func (q rangeQuery) Resolve(r *Resolver) {
	r.Offset = q[0]
	r.Limit = q[1]
}

type orderQuery uint8

// Apply implements the Query interface
func (q orderQuery) Apply(p *post.Post) {}

// Resolve implements the Query interface
func (q orderQuery) Resolve(r *Resolver) {
	r.Order = uint8(q)
}

type filterQuery struct {
	T func(post.Post) bool
}

// Apply implements the Query interface
func (f filterQuery) Apply(p *post.Post) {}

// Resolve implements the Query interface
func (f filterQuery) Resolve(r *Resolver) {
	r.Filters = append(r.Filters, f.T)
}
