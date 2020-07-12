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
	switch query.QueryType() {
	case "id":
		q, ok := query.(idQuery)
		if ok {
			if len(q) == 20 {
				id, err := xid.FromString(string(q))
				if err == nil {
					res.ID = id.Bytes()
				}
			}
		}

	case "slug":
		q, ok := query.(slugQuery)
		if ok {
			res.Slug = toLowerTrim(string(q))
		}

	case "type":
		q, ok := query.(typeQuery)
		if ok {
			res.Type = toLowerTrim(string(q))
		}

	case "owner":
		q, ok := query.(ownerQuery)
		if ok {
			res.Owner = toLowerTrim(string(q))
		}

	case "status":
		q, ok := query.(statusQuery)
		if ok {
			res.Status = uint8(q)
		}

	case "tags":
		q, ok := query.(tagsQuery)
		if ok {
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
		}

	case "keys":
		q, ok := query.(keysQuery)
		if ok {
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
		}

	case "range":
		q, ok := query.(rangeQuery)
		if ok {
			if len(q.after) == 20 {
				id, err := xid.FromString(q.after)
				if err == nil {
					res.Aftar = id.Bytes()
				}
			}
			res.Limit = q.limit
		}

	case "order":
		q, ok := query.(orderQuery)
		if ok {
			res.Order = uint8(q)
		}
	}
}
