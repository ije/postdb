package q

import (
	"testing"

	"github.com/ije/postdb/internal/post"
)

func TestQuery(t *testing.T) {
	var res Resolver
	id := post.NewID()
	id2 := post.NewID()
	for _, q := range []Query{
		ID(id),
		ID(id2),
		Owner("admin"),
		Status(1),
		Tags("hello", "world"),
		Tags("world", "世界"),
		Select("title", "date ", "content", "content"),
		Anchor(id2),
		Offset(2),
		Limit(100),
		Order(DESC),
	} {
		q.Resolve(&res)
	}

	toBe(t, "IDs", len(res.IDs), 2)
	toBe(t, "Owner", res.Owner, "admin")
	toBe(t, "Tags", len(res.Tags), 3)
	toBe(t, "Keys", len(res.Keys), 3)
	toBe(t, "Anchor", res.Anchor, id2)
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
