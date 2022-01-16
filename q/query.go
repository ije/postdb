package q

import (
	"github.com/ije/postdb/post"
	"github.com/ije/postdb/util"
)

const (
	// DESC specifies the order of DESC
	DESC uint8 = iota
	// ASC specifies the order of ASC
	ASC
	// todo:
	// RANK_DESC specifies the order of DESC by Rank
	// RANK_DESC
	// RANK_ASC specifies the order of ASC by Rank
	// RANK_ASC
)

// A Query inferface
type Query interface {
	Resolve(*Resolver)
	Apply(*post.Post)
}

// ID returns an id Query
func ID(id string) Query {
	return idQuery(id)
}

// IDs returns a IDs Query
func IDs(ids ...string) Query {
	return idsQuery(util.NoRepeat(ids))
}

// Alias returns a alias Query
func Alias(alias string) Query {
	return aliasQuery(alias)
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
	return tagsQuery(util.NoRepeat(tags))
}

// Select returns a keys Query
func Select(keys ...string) Query {
	return keysQuery(util.NoRepeat(keys))
}

// Anchor returns a anchor Query
func Anchor(id string) Query {
	return anchorQuery(id)
}

// Offset returns an offset Query
func Offset(id uint32) Query {
	return offsetQuery(id)
}

// Limit returns a limit Query
func Limit(limit uint32) Query {
	return limitQuery(limit)
}

// Range returns a range Query
func Range(offset uint32, limit uint32) Query {
	return rangeQuery([2]uint32{offset, limit})
}

// Order returns a order Query
func Order(order uint8) Query {
	return orderQuery(order)
}

// Filter returns a filter Query
func Filter(fn func(post.Post) bool) Query {
	return filterQuery{fn}
}
