package q

import (
	"encoding/binary"

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
type rangeQuery [17]byte
type limitQuery uint32
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
	a := make(tagsQuery, len(tags))
	i := 0
	for _, s := range tags {
		tag := toLowerTrim(s)
		if tag != "" {
			a[i] = tag
			i++
		}
	}
	return a[:i]
}

// Keys returns a keys Query
func Keys(keys ...string) Query {
	a := make(keysQuery, len(keys))
	i := 0
	for _, s := range keys {
		tag := toLowerTrim(s)
		if tag != "" {
			a[i] = tag
			i++
		}
	}
	return a[:i]
}

// Range returns a range Query
func Range(after string, limit uint32) Query {
	var q rangeQuery
	if len(after) == 20 {
		q[0] = 1
		id, err := xid.FromString(after)
		if err == nil {
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
func (order orderQuery) QueryType() string {
	return "order"
}
