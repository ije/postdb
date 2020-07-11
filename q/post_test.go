package q

import (
	"strings"
	"testing"
)

func TestPostMeta(t *testing.T) {
	post := NewPost("blog")
	for _, q := range []Query{
		Slug("Hello World"),
		ACL(PUBLIC_READ_WRITE),
		Owner("admin"),
		Tags("hello", "world"),
		KV{
			"title": []byte("Hello World!"),
			"date":  []byte("2020-01-01"),
		},
	} {
		post.ApplyQuery(q)
	}

	metadata := post.MetaData()
	_post, err := ParsePostMeta(metadata)
	if err != nil {
		t.Fatal(err)
	}

	toBe(t, "_post.ID", string(_post.ID), string(post.ID))
	toBe(t, "_post.Type", _post.Type, post.Type)
	toBe(t, "_post.ACL", _post.ACL, post.ACL)
	toBe(t, "_post.Owner", _post.Owner, post.Owner)
	toBe(t, "_post.Crtime", _post.Crtime, post.Crtime)
	toBe(t, "_post.Tags", strings.Join(_post.Tags, ","), strings.Join(post.Tags, ","))
	toBe(t, "_post.KV.title", string(_post.KV["title"]), string(post.KV["title"]))
	toBe(t, "_post.KV.date", string(_post.KV["date"]), string(post.KV["date"]))

	t.Log(_post)
}

func toBe(t *testing.T, name string, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("%s(%v) should equal to %v", name, a, b)
	}
}
