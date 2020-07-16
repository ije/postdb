package q

import (
	"encoding/binary"
)

// Resolver to resolves query
type Resolver struct {
	ID     []byte
	Slug   string
	Type   string
	Owner  string
	Status uint8
	Tags   []string
	Keys   []string
	After  []byte
	KV     KV
	Limit  uint32
	Order  uint8
}

// Apply applies a query
func (res *Resolver) Apply(query Query) {
	switch q := query.(type) {
	case ObjectID:
		if !q.IsNil() {
			res.ID = q.Bytes()
		}

	case slugQuery:
		res.Slug = toLowerTrim(string(q))

	case typeQuery:
		res.Type = toLowerTrim(string(q))

	case ownerQuery:
		res.Owner = toLowerTrim(string(q))

	case statusQuery:
		res.Status = uint8(q)

	case tagsQuery:
		l := len(q)
		if l > 0 {
			n := len(res.Tags)
			tags := make([]string, n+l)
			if n > 0 {
				copy(tags, res.Tags)
			}
			copy(tags[n:], q)
			res.Tags = tags
		}

	case keysQuery:
		l := len(q)
		if l > 0 {
			n := len(res.Keys)
			keys := make([]string, n+l)
			if n > 0 {
				copy(keys, res.Keys)
			}
			copy(keys[n:], q)
			res.Keys = keys
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
