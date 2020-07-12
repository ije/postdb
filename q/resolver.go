package q

import (
	"github.com/rs/xid"
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
	Aftar  []byte
	Limit  int
	Order  uint8
}

// Apply applies a query
func (res *Resolver) Apply(query Query) {
	switch q := query.(type) {
	case idQuery:
		if len(q) == 20 {
			id, err := xid.FromString(string(q))
			if err == nil {
				res.ID = id.Bytes()
			}
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

	case rangeQuery:
		if len(q.after) == 20 {
			id, err := xid.FromString(q.after)
			if err == nil {
				res.Aftar = id.Bytes()
			}
		}
		res.Limit = q.limit

	case orderQuery:
		res.Order = uint8(q)
	}
}
