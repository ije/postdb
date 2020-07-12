package q

// A Query inferface
type Query interface {
	QueryType() string
}

// A KV map
type KV map[string][]byte

type idQuery string
type slugQuery string
type typeQuery string
type ownerQuery string
type statusQuery uint8
type tagsQuery []string
type keysQuery []string
type rangeQuery struct {
	after string
	limit int
}
type orderQuery uint8

// ID returns a id Query
func ID(id string) Query {
	return idQuery(id)
}

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
func Tags(tag string, extra ...string) Query {
	tags := make(tagsQuery, 1+len(extra))
	i := 0
	t := toLowerTrim(tag)
	if t != "" {
		tags[i] = t
		i++
	}
	for _, s := range extra {
		t := toLowerTrim(s)
		if t != "" {
			tags[i] = t
			i++
		}
	}
	return tags[:i]
}

// Keys returns a keys Query
func Keys(key string, extra ...string) Query {
	keys := make(keysQuery, 1+len(extra))
	i := 0
	if key != "" {
		keys[i] = key
		i++
	}
	for _, k := range extra {
		if k != "" {
			keys[i] = k
			i++
		}
	}
	return keys[:i]
}

// Range returns a range Query
func Range(after string, limit int) Query {
	return rangeQuery{
		after: after,
		limit: limit,
	}
}

// Order returns a order Query
func Order(order uint8) Query {
	return orderQuery(order)
}

// QueryType implements the Query interface
func (kv KV) QueryType() string {
	return "kv"
}

// QueryType implements the Query interface
func (q idQuery) QueryType() string {
	return "id"
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
