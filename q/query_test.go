package q

import (
	"testing"

	"github.com/rs/xid"
)

func TestQuery(t *testing.T) {
	var res Resolver
	for _, q := range []Query{
		ID("id1"),
		ID(xid.New().String()),
		Owner(7),
		Status(1),
		Tags("hello", "world"),
		Tags("world", "世界"),
		Select("title", "date ", "content", "content"),
		Offset(2),
		Limit(100),
		Order(DESC),
	} {
		q.Resolve(&res)
	}

	toBe(t, "IDs", len(res.IDs), 1)
	toBe(t, "Owner", res.Owner, uint32(7))
	toBe(t, "Tags", len(res.Tags), 3)
	toBe(t, "Keys", len(res.Keys), 3)
	toBe(t, "Offset", res.Offset, uint32(2))
	toBe(t, "Limit", res.Limit, uint32(100))
	toBe(t, "Order", res.Order, DESC)

	t.Log(res)
}

func toBe(t *testing.T, name string, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("the %s should equal to %v, but %v", name, b, a)
	}
}
