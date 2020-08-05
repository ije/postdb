package q

import (
	"testing"
)

func TestQuery(t *testing.T) {
	var res Resolver
	id := NewID()
	aid := NewID()
	for _, q := range []Query{
		ID(id),
		ID(aid),
		Owner("admin"),
		Status(1),
		Tags("hello", "world"),
		Tags("world", "世界"),
		K("title", "content", "*", "content"),
		Offset(aid),
		Anchor(aid),
		Limit(100),
		Order(DESC),
	} {
		q.Resolve(&res)
	}

	toBe(t, "IDs", len(res.IDs), 2)
	toBe(t, "Owner", res.Owner, "admin")
	toBe(t, "Tags", len(res.Tags), 3)
	toBe(t, "Keys", len(res.Keys), 3)
	toBe(t, "Anchor", res.Anchor, aid)
	toBe(t, "Offset", res.Offset, aid)
	toBe(t, "Limit", res.Limit, uint32(100))
	toBe(t, "Order", res.Order, DESC)

	t.Log(res)
}
