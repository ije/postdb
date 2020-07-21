package q

import (
	"encoding/binary"
	"strings"
)

// A Setter inferface
// type Setter interface {
// 	Set(*Post)
// 	Error() error
// }

// A Query inferface
// type Query interface {
// 	Apply(*Resolver)
// 	Error() error
// }

// A Query inferface
type Query interface {
	Error() error
}

type idsQuery [][]byte
type aliasQuery string
type ownerQuery string
type statusQuery uint8
type tagsQuery []string
type keysQuery []string
type afterQuery [13]byte
type limitQuery uint32
type rangeQuery [17]byte
type orderQuery uint8

// IDs returns a IDs Query
func IDs(ids ...string) Query {
	set := map[string]struct{}{}
	a := make([][]byte, len(ids))
	i := 0
	for _, id := range ids {
		xid := ID(id)
		if !xid.IsNil() {
			_, ok := set[id]
			if !ok {
				set[id] = struct{}{}
				a[i] = xid.Bytes()
				i++
			}
		}
	}
	return idsQuery(a[:i])
}

// Alias returns a alias Query
func Alias(alias string) Query {
	return aliasQuery(strings.ReplaceAll(strings.ToLower(strings.TrimSpace(alias)), " ", "-"))
}

// Owner returns a owner Query
func Owner(name string) Query {
	return ownerQuery(strings.TrimSpace(name))
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
	var q afterQuery
	xid := ID(id)
	if !xid.IsNil() {
		q[0] = 1
		copy(q[1:], xid.Bytes())
	}
	return q
}

// Limit returns a limit Query
func Limit(limit uint8) Query {
	return limitQuery(limit)
}

// Range returns a range Query.
//
// `postdb.GetPosts(q.Range("bs7pobh8d3b21ducpaqg", 100))` equals `postdb.GetPosts(q.After("bs7pobh8d3b21ducpaqg"), q.Limit(100))`
func Range(after string, limit uint32) Query {
	var q rangeQuery
	if len(after) == 20 && limit > 0 {
		xid := ID(after)
		if !xid.IsNil() {
			q[0] = 1
			copy(q[1:], xid.Bytes())
			binary.BigEndian.PutUint32(q[13:], limit)
		}
	}
	return q
}

// Order returns a order Query
func Order(order uint8) Query {
	return orderQuery(order)
}

// Error implements the Query interface
func (q idsQuery) Error() error {
	return nil
}

// Error implements the Query interface
func (q aliasQuery) Error() error {
	return nil
}

// Error implements the Query interface
func (q ownerQuery) Error() error {
	return nil
}

// Error implements the Query interface
func (status statusQuery) Error() error {
	return nil
}

// Error implements the Query interface
func (q tagsQuery) Error() error {
	return nil
}

// Error implements the Query interface
func (q keysQuery) Error() error {
	return nil
}

// Error implements the Query interface
func (q rangeQuery) Error() error {
	return nil
}

// Error implements the Query interface
func (q afterQuery) Error() error {
	return nil
}

// Error implements the Query interface
func (q limitQuery) Error() error {
	return nil
}

// Error implements the Query interface
func (order orderQuery) Error() error {
	return nil
}
