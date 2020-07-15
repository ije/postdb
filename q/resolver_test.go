package q

import (
	"strings"
	"testing"
)

func TestQuery(t *testing.T) {
	var res Resolver
	for _, q := range []Query{
		ID(NewID().String()),
		Slug("hello-world"),
		Type("news"),
		Owner("admin"),
		Status(123),
		Keys("title", "content"),
		KV{"title": []byte("Hello World")},
		Tags("hello", "world"),
		Range(NewID().String(), 1024),
		Order(DESC),
	} {
		res.Apply(q)
	}

	toBe(t, "ID", len(res.ID), 12)
	toBe(t, "Slug", res.Slug, "hello-world")
	toBe(t, "Type", res.Type, "news")
	toBe(t, "Owner", res.Owner, "admin")
	toBe(t, "Status", res.Status, uint8(123))
	toBe(t, "Tags", strings.Join(res.Tags, " "), "hello world")
	toBe(t, "Keys", strings.Join(res.Keys, " "), "title content")
	toBe(t, "Keys.title", string(res.KV["title"]), "Hello World")
	toBe(t, "Aftar", len(res.After), 12)
	toBe(t, "Limit", res.Limit, 1024)
	toBe(t, "Order", res.Order, DESC)

	t.Log(res)
}
