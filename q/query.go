package q

import (
	"github.com/rs/xid"
)

type Query interface {
	QueryType() string
}

type QueryResult struct {
	Type  string
	Owner string
	Tags  []string
	Keys  []string
	Aftar []byte
	Limit int
	Order Order
}

// ApplyQuery applies query
func (ret *QueryResult) ApplyQuery(query Query) {
	switch query.QueryType() {
	case "type":
		q, ok := query.(typeQuery)
		if ok {
			ret.Type = toLowerTrim(string(q))
		}

	case "owner":
		q, ok := query.(ownerQuery)
		if ok {
			ret.Owner = toLowerTrim(string(q))
		}

	case "tags":
		q, ok := query.(tagsQuery)
		if ok {
			l := len(q)
			if l > 0 {
				n := len(ret.Tags)
				tags := make([]string, n+l)
				if n > 0 {
					copy(tags, ret.Tags)
				}
				copy(tags[n:], q)
				ret.Tags = tags
			}
		}

	case "keys":
		q, ok := query.(Keys)
		if ok {
			l := len(q)
			if l > 0 {
				n := len(ret.Keys)
				keys := make([]string, n+l)
				if n > 0 {
					copy(keys, ret.Keys)
				}
				copy(keys[n:], q)
				ret.Keys = keys
			}
		}

	case "range":
		q, ok := query.(rangeQuery)
		if ok {
			if len(q.after) == 20 {
				id, err := xid.FromString(q.after)
				if err == nil {
					ret.Aftar = id.Bytes()
				}
			}
			ret.Limit = q.limit
		}

	case "order":
		q, ok := query.(Order)
		if ok {
			ret.Order = q
		}
	}

}
