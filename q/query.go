package q

import (
	"strings"
)

// Resolver to save query resolves
type Resolver struct {
	ID              string
	IDs             []string
	Alias           string
	Owner           string
	Status          uint8
	HasStatus       bool
	Tags            map[string]struct{}
	Keys            map[string]struct{}
	KeysHasWildcard bool
	After           string
	Limit           uint32
	Order           uint8
}

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
type afterQuery string
type limitQuery uint32
type orderQuery uint8

// ID returns an id Query
func ID(id string) Query {
	return idQuery(id)
}

// IDs returns a IDs Query
func IDs(ids ...string) Query {
	set := map[string]struct{}{}
	a := make([]string, len(ids))
	i := 0
	for _, id := range ids {
		if id != "" {
			_, ok := set[id]
			if !ok {
				set[id] = struct{}{}
				a[i] = id
				i++
			}
		}
	}
	return idsQuery(a[:i])
}

// Alias returns a alias Query
func Alias(alias string) Query {
	p := strings.Split(alias, " ")
	a := make([]string, len(p))
	i := 0
	for _, s := range p {
		s = strings.ToLower(strings.TrimSpace(alias))
		if s != "" {
			a[i] = s
			i++
		}
	}
	return aliasQuery(strings.Join(a[:i], "-"))
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
	set := map[string]struct{}{}
	a := make([]string, len(tags))
	i := 0
	for _, s := range tags {
		tag := strings.ToLower(strings.TrimSpace(s))
		if tag != "" {
			_, ok := set[tag]
			if !ok {
				set[tag] = struct{}{}
				a[i] = tag
				i++
			}
		}
	}
	return tagsQuery(a[:i])
}

// Keys returns a keys Query
func Keys(keys ...string) Query {
	set := map[string]struct{}{}
	a := make([]string, len(keys))
	i := 0
	for _, s := range keys {
		key := strings.TrimSpace(s)
		if key != "" {
			_, ok := set[key]
			if !ok {
				set[key] = struct{}{}
				a[i] = key
				i++
			}
		}
	}
	return keysQuery(a[:i])
}

// K is a shortcut for Key
func K(keys ...string) Query {
	return Keys(keys...)
}

// After returns a after Query
func After(id string) Query {
	return afterQuery(id)
}

// Limit returns a limit Query
func Limit(limit uint8) Query {
	return limitQuery(limit)
}

// Order returns a order Query
func Order(order uint8) Query {
	return orderQuery(order)
}

// Apply implements the Query interface
func (q idQuery) Apply(p *Post) {}

// Resolve implements the Query interface
func (q idQuery) Resolve(r *Resolver) {
	r.ID = string(q)
}

// Apply implements the Query interface
func (q idsQuery) Apply(p *Post) {}

// Resolve implements the Query interface
func (q idsQuery) Resolve(r *Resolver) {
	r.IDs = q
}

// Apply implements the Query interface
func (q aliasQuery) Apply(p *Post) {
	p.Alias = string(q)
}

// Resolve implements the Query interface
func (q aliasQuery) Resolve(r *Resolver) {
	r.Alias = string(q)
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
	r.Status = uint8(q)
	r.HasStatus = true
}

// Apply implements the Query interface
func (q tagsQuery) Apply(p *Post) {
	p.Tags = q
}

// Resolve implements the Query interface
func (q tagsQuery) Resolve(r *Resolver) {
	if r.Tags == nil {
		r.Tags = map[string]struct{}{}
	}
	for _, tag := range q {
		r.Tags[tag] = struct{}{}
	}
}

// Apply implements the Query interface
func (q keysQuery) Apply(p *Post) {}

// Resolve implements the Query interface
func (q keysQuery) Resolve(r *Resolver) {
	if r.Keys == nil {
		r.Keys = map[string]struct{}{}
	}
	for _, s := range q {
		if s == "*" {
			r.KeysHasWildcard = true
		}
		r.Keys[s] = struct{}{}
	}
}

// Apply implements the Query interface
func (q afterQuery) Apply(p *Post) {}

// Resolve implements the Query interface
func (q afterQuery) Resolve(r *Resolver) {
	r.After = string(q)
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
