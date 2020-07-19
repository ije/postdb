package q

import (
	"testing"
)

func TestQuery(t *testing.T) {
	var res Resolver
	id := NewID()
	afterID := NewID()
	for _, q := range []Query{
		ID(id.String()),
		Slug("hello-world"),
		Type("news"),
		Owner("admin"),
		Status(123),
		Tags("hello", "world"),
		Tags("world", "世界"),
		Keys("title", "content", "*", "content"),
		KV{"title": []byte("Hello World")},
		After(afterID.String()),
		Limit(100),
		Order(DESC),
	} {
		res.Apply(q)
	}

	toBe(t, "ID", string(res.ID), string(id.Bytes()))
	toBe(t, "Slug", res.Slug, "hello-world")
	toBe(t, "Type", res.Type, "news")
	toBe(t, "Owner", res.Owner, "admin")
	toBe(t, "Status", res.Status, uint8(123))
	toBe(t, "Tags", len(res.Tags), 3)
	toBe(t, "KVKeys", len(res.KVKeys), 3)
	toBe(t, "KVWildcard", res.KVWildcard, true)
	toBe(t, "KV.title", string(res.KV["title"]), "Hello World")
	toBe(t, "Aftar", string(res.After), string(afterID.Bytes()))
	toBe(t, "Limit", res.Limit, uint32(100))
	toBe(t, "Order", res.Order, DESC)

	t.Log(res)
}
