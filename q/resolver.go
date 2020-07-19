package q

import (
	"encoding/binary"
)

// Resolver to resolves query
type Resolver struct {
	ID          []byte
	Slug        string
	Type        string
	Owner       string
	Status      uint8
	Tags        []string
	Keys        []string
	WildcardKey bool
	KV          KV
	After       []byte
	Limit       uint32
	Order       uint8
}

// Apply applies a query
func (res *Resolver) Apply(query Query) {
	switch q := query.(type) {
	case ObjectID:
		if !q.IsNil() {
			res.ID = q.Bytes()
		}

	case slugQuery:
		res.Slug = string(q)

	case typeQuery:
		res.Type = string(q)

	case ownerQuery:
		res.Owner = string(q)

	case statusQuery:
		res.Status = uint8(q)

	case tagsQuery:
		l := len(q)
		if l > 0 {
			set := map[string]struct{}{}
			n := len(res.Tags)
			a := make([]string, n+l)
			for i, s := range res.Tags {
				set[s] = struct{}{}
				a[i] = s
			}
			for _, s := range q {
				_, ok := set[s]
				if !ok {
					set[s] = struct{}{}
					a[n] = s
					n++
				}
			}
			res.Tags = a[:n]
		}

	case keysQuery:
		l := len(q)
		if l > 0 {
			set := map[string]struct{}{}
			n := len(res.Keys)
			a := make([]string, n+l)
			for i, s := range res.Keys {
				set[s] = struct{}{}
				a[i] = s
			}
			for _, s := range q {
				if s == "*" {
					res.WildcardKey = true
				}
				_, ok := set[s]
				if !ok {
					set[s] = struct{}{}
					a[n] = s
					n++
				}
			}
			res.Keys = a[:n]
		}

	case KV:
		if res.KV == nil {
			res.KV = KV{}
		}
		for k, v := range q {
			if len(k) > 0 {
				res.KV[k] = v
			}
		}

	case rangeQuery:
		if q[0] == 1 {
			res.After = q[1:13]
			res.Limit = binary.BigEndian.Uint32(q[13:])
		}

	case afterQuery:
		if q[0] == 1 {
			res.After = q[1:]
		}

	case limitQuery:
		res.Limit = uint32(q)

	case orderQuery:
		res.Order = uint8(q)
	}
}
