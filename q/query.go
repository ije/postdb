package q

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
	Apply(*Post)
}

type idQuery string
type idsQuery []string
type aliasQuery string
type ownerQuery string
type statusQuery uint8
type tagsQuery []string
type keysQuery []string
type anchorQuery string
type offsetQuery string
type limitQuery uint32
type orderQuery uint8
type filterQuery struct {
	T func(Post) bool
}

// ID returns an id Query
func ID(id string) Query {
	return idQuery(id)
}

// IDs returns a IDs Query
func IDs(ids ...string) Query {
	return idsQuery(StringSet(ids))
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
	return tagsQuery(StringSet(tags))
}

// Keys returns a keys Query
func Keys(keys ...string) Query {
	return keysQuery(StringSet(keys))
}

// K is a shortcut for Key
func K(keys ...string) Query {
	return Keys(keys...)
}

// Anchor returns a anchor Query
func Anchor(id string) Query {
	return anchorQuery(id)
}

// Offset returns an offset Query
func Offset(id string) Query {
	return offsetQuery(id)
}

// Limit returns a limit Query
func Limit(limit uint32) Query {
	return limitQuery(limit)
}

// Order returns a order Query
func Order(order uint8) Query {
	return orderQuery(order)
}

// Filter returns a filter Query
func Filter(fn func(Post) bool) Query {
	return filterQuery{fn}
}

// Apply implements the Query interface
func (q idQuery) Apply(p *Post) {}

// Resolve implements the Query interface
func (q idQuery) Resolve(r *Resolver) {
	r.IDs = append(r.IDs, string(q))
}

// Apply implements the Query interface
func (q idsQuery) Apply(p *Post) {}

// Resolve implements the Query interface
func (q idsQuery) Resolve(r *Resolver) {
	r.IDs = append(r.IDs, q...)
}

// Apply implements the Query interface
func (q aliasQuery) Apply(p *Post) {
	p.Alias = string(q)
}

// Resolve implements the Query interface
func (q aliasQuery) Resolve(r *Resolver) {
	r.IDs = append(r.IDs, string(q))
}

// Apply implements the Query interface
func (q ownerQuery) Apply(p *Post) {
	p.Owner = string(q)
}

// Resolve implements the Query interface
func (q ownerQuery) Resolve(r *Resolver) {
	r.Owner = string(q)
}

// Apply implements the Query interface
func (q statusQuery) Apply(p *Post) {
	p.Status = uint8(q)
}

// Resolve implements the Query interface
func (q statusQuery) Resolve(r *Resolver) {
	r.Filters = append(r.Filters, func(p Post) bool {
		return p.Status == uint8(q)
	})
}

// Apply implements the Query interface
func (q tagsQuery) Apply(p *Post) {
	p.Tags = q
}

// Resolve implements the Query interface
func (q tagsQuery) Resolve(r *Resolver) {
	tags := append(r.Tags, q...)
	r.Tags = StringSet(tags)
}

// Apply implements the Query interface
func (q keysQuery) Apply(p *Post) {}

// Resolve implements the Query interface
func (q keysQuery) Resolve(r *Resolver) {
	keys := append(r.Keys, q...)
	r.Keys = StringSet(keys)
}

// Apply implements the Query interface
func (q anchorQuery) Apply(p *Post) {}

// Resolve implements the Query interface
func (q anchorQuery) Resolve(r *Resolver) {
	r.Anchor = string(q)
}

// Apply implements the Query interface
func (q offsetQuery) Apply(p *Post) {}

// Resolve implements the Query interface
func (q offsetQuery) Resolve(r *Resolver) {
	r.Offset = string(q)
}

// Apply implements the Query interface
func (q limitQuery) Apply(p *Post) {}

// Resolve implements the Query interface
func (q limitQuery) Resolve(r *Resolver) {
	r.Limit = uint32(q)
}

// Apply implements the Query interface
func (q orderQuery) Apply(p *Post) {}

// Resolve implements the Query interface
func (q orderQuery) Resolve(r *Resolver) {
	r.Order = uint8(q)
}

// Apply implements the Query interface
func (f filterQuery) Apply(p *Post) {}

// Resolve implements the Query interface
func (f filterQuery) Resolve(r *Resolver) {
	r.Filters = append(r.Filters, f.T)
}
