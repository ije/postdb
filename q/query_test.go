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
		IDs(id, aid),
		Owner("admin"),
		Status(123),
		Tags("hello", "world"),
		Tags("world", "世界"),
		K("title", "content", "*", "content"),
		Offset(aid),
		Limit(100),
		Order(DESC),
	} {
		q.Resolve(&res)
	}

	toBe(t, "ID", res.ID, id)
	toBe(t, "IDs", len(res.IDs), 2)
	toBe(t, "Owner", res.Owner, "admin")
	toBe(t, "Status", res.Status, uint8(123))
	toBe(t, "Tags", len(res.Tags), 3)
	toBe(t, "Keys", len(res.Keys), 3)
	toBe(t, "HasWildcardKey", res.HasWildcardKey, true)
	toBe(t, "Offset", res.Offset, aid)
	toBe(t, "Limit", res.Limit, uint32(100))
	toBe(t, "Order", res.Order, DESC)

	t.Log(res)
}
