package q

import (
	"strings"
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
		Keys("title", "content", "*"),
		KV{"title": []byte("Hello World")},
		Tags("hello", "world"),
		Tags("world", "世界"),
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
	toBe(t, "Tags", strings.Join(res.Tags, " "), "hello world 世界")
	toBe(t, "Keys", strings.Join(res.Keys, " "), "title content *")
	toBe(t, "WildcardKey", res.WildcardKey, true)
	toBe(t, "Keys.title", string(res.KV["title"]), "Hello World")
	toBe(t, "Aftar", string(res.After), string(afterID.Bytes()))
	toBe(t, "Limit", res.Limit, uint32(100))
	toBe(t, "Order", res.Order, DESC)

	t.Log(res)
}
