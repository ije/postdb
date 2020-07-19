package q

import (
	"encoding/binary"
)

// Resolver to resolves query
type Resolver struct {
	ID         []byte
	Slug       string
	Type       string
	Owner      string
	Status     uint8
	HasStatus  bool
	Tags       map[string]struct{}
	KV         KV
	KVKeys     map[string]struct{}
	KVWildcard bool
	After      []byte
	Limit      uint32
	Order      uint8
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
		res.HasStatus = true

	case tagsQuery:
		if res.Tags == nil {
			res.Tags = map[string]struct{}{}
		}
		for _, s := range q {
			res.Tags[s] = struct{}{}
		}

	case keysQuery:
		if res.KVKeys == nil {
			res.KVKeys = map[string]struct{}{}
		}
		for _, s := range q {
			if s == "*" {
				res.KVWildcard = true
			}
			res.KVKeys[s] = struct{}{}
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
