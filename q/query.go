package q

import (
	"encoding/binary"
	"strings"

	"github.com/rs/xid"
)

// A Query inferface
type Query interface {
	QueryType() string
}

type slugQuery string
type typeQuery string
type ownerQuery string
type statusQuery uint8
type tagsQuery []string
type keysQuery []string
type afterQuery [13]byte
type limitQuery uint32
type rangeQuery [17]byte
type orderQuery uint8

// Slug returns a slug Query
func Slug(slug string) Query {
	return slugQuery(slug)
}

// Type returns a type Query
func Type(t string) Query {
	return typeQuery(t)
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
		tag := toLowerTrim(s)
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

// After returns a after Query
func After(id string) Query {
	var q afterQuery
	if len(id) == 20 {
		xid, err := xid.FromString(id)
		if err == nil {
			q[0] = 1
			copy(q[1:], xid.Bytes())
		}
	}
	return q
}

// Limit returns a limit Query
func Limit(limit uint8) Query {
	return limitQuery(limit)
}

// Range returns a range Query.
//
// `postdb.GetPosts(q.Range("bs7pobh8d3b21ducpaqg", 100))`  equals `postdb.GetPosts(q.After("bs7pobh8d3b21ducpaqg"), q.Limit(100))`
func Range(after string, limit uint32) Query {
	var q rangeQuery
	if len(after) == 20 && limit > 0 {
		id, err := xid.FromString(after)
		if err == nil {
			q[0] = 1
			copy(q[1:], id.Bytes())
			binary.BigEndian.PutUint32(q[13:], limit)
		}
	}
	return q
}

// Order returns a order Query
func Order(order uint8) Query {
	return orderQuery(order)
}

// QueryType implements the Query interface
func (q slugQuery) QueryType() string {
	return "slug"
}

// QueryType implements the Query interface
func (q typeQuery) QueryType() string {
	return "type"
}

// QueryType implements the Query interface
func (q ownerQuery) QueryType() string {
	return "owner"
}

// QueryType implements the Query interface
func (status statusQuery) QueryType() string {
	return "status"
}

// QueryType implements the Query interface
func (q tagsQuery) QueryType() string {
	return "tags"
}

// QueryType implements the Query interface
func (q keysQuery) QueryType() string {
	return "keys"
}

// QueryType implements the Query interface
func (q rangeQuery) QueryType() string {
	return "range"
}

// QueryType implements the Query interface
func (q afterQuery) QueryType() string {
	return "after"
}

// QueryType implements the Query interface
func (q limitQuery) QueryType() string {
	return "limit"
}

// QueryType implements the Query interface
func (order orderQuery) QueryType() string {
	return "order"
}
