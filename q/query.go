package q

import "github.com/postui/postdb/post"

const (
	// DESC specifies the order of DESC
	DESC uint8 = iota
	// ASC specifies the order of ASC
	ASC
	// todo:
	// RANK_DESC specifies the order of DESC by Rank
	// RANK_DESC
	// RANK_ASC specifies the order of ASC by Rank
	// RANK_ASC
)

// A Query inferface
type Query interface {
	Resolve(*Resolver)
	Apply(*post.Post)
}

type idQuery string
type idsQuery []string
type aliasQuery string
type ownerQuery string
type statusQuery uint8
type tagsQuery []string
type keysQuery []string
type anchorQuery string
type offsetQuery uint32
type limitQuery uint32
type rangeQuery [2]uint32
type orderQuery uint8
type filterQuery struct {
	T func(post.Post) bool
}

// ID returns an id Query
func ID(id string) Query {
	return idQuery(id)
}

// IDs returns a IDs Query
func IDs(ids ...string) Query {
	return idsQuery(noRepeat(ids))
}

// Alias returns a alias Query
func Alias(alias string) Query {
	return aliasQuery(alias)
}

// Owner returns a owner Query
func Owner(name string) Query {
	return ownerQuery(name)
}

// Status returns a status Query
func Status(status uint8) Query {
	return statusQuery(status)
}

// Tags returns a tags Query
func Tags(tags ...string) Query {
	return tagsQuery(noRepeat(tags))
}

// Select returns a keys Query
func Select(keys ...string) Query {
	return keysQuery(noRepeat(keys))
}

// Anchor returns a anchor Query
func Anchor(id string) Query {
	return anchorQuery(id)
}

// Offset returns an offset Query
func Offset(id uint32) Query {
	return offsetQuery(id)
}

// Limit returns a limit Query
func Limit(limit uint32) Query {
	return limitQuery(limit)
}

// Range returns a range Query
func Range(offset uint32, limit uint32) Query {
	return rangeQuery([2]uint32{offset, limit})
}

// Order returns a order Query
func Order(order uint8) Query {
	return orderQuery(order)
}

// Filter returns a filter Query
func Filter(fn func(post.Post) bool) Query {
	return filterQuery{fn}
}

// Apply implements the Query interface
func (q idQuery) Apply(p *post.Post) {}

// Resolve implements the Query interface
func (q idQuery) Resolve(r *Resolver) {
	r.IDs = append(r.IDs, string(q))
}

// Apply implements the Query interface
func (q idsQuery) Apply(p *post.Post) {}

// Resolve implements the Query interface
func (q idsQuery) Resolve(r *Resolver) {
	r.IDs = append(r.IDs, q...)
}

// Apply implements the Query interface
func (q aliasQuery) Apply(p *post.Post) {
	p.Alias = string(q)
}

// Resolve implements the Query interface
func (q aliasQuery) Resolve(r *Resolver) {
	r.IDs = append(r.IDs, string(q))
}

// Apply implements the Query interface
func (q ownerQuery) Apply(p *post.Post) {
	p.Owner = string(q)
}

// Resolve implements the Query interface
func (q ownerQuery) Resolve(r *Resolver) {
	r.Owner = string(q)
}

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

// Apply implements the Query interface
func (q tagsQuery) Apply(p *post.Post) {
	p.Tags = q
}

// Resolve implements the Query interface
func (q tagsQuery) Resolve(r *Resolver) {
	tags := append(r.Tags, q...)
	r.Tags = noRepeat(tags)
}

// Apply implements the Query interface
func (q keysQuery) Apply(p *post.Post) {}

// Resolve implements the Query interface
func (q keysQuery) Resolve(r *Resolver) {
	keys := append(r.Keys, q...)
	r.Keys = noRepeat(keys)
}

// Apply implements the Query interface
func (q anchorQuery) Apply(p *post.Post) {}

// Resolve implements the Query interface
func (q anchorQuery) Resolve(r *Resolver) {
	r.Anchor = string(q)
}

// Apply implements the Query interface
func (q offsetQuery) Apply(p *post.Post) {}

// Resolve implements the Query interface
func (q offsetQuery) Resolve(r *Resolver) {
	r.Offset = uint32(q)
}

// Apply implements the Query interface
func (q limitQuery) Apply(p *post.Post) {}

// Resolve implements the Query interface
func (q limitQuery) Resolve(r *Resolver) {
	r.Limit = uint32(q)
}

// Apply implements the Query interface
func (q rangeQuery) Apply(p *post.Post) {}

// Resolve implements the Query interface
func (q rangeQuery) Resolve(r *Resolver) {
	r.Offset = q[0]
	r.Limit = q[1]
}

// Apply implements the Query interface
func (q orderQuery) Apply(p *post.Post) {}

// Resolve implements the Query interface
func (q orderQuery) Resolve(r *Resolver) {
	r.Order = uint8(q)
}

// Apply implements the Query interface
func (f filterQuery) Apply(p *post.Post) {}

// Resolve implements the Query interface
func (f filterQuery) Resolve(r *Resolver) {
	r.Filters = append(r.Filters, f.T)
}

func noRepeat(arr []string) []string {
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
