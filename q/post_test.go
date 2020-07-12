package q

import (
	"strings"
	"testing"
)

func TestPostMeta(t *testing.T) {
	post := NewPost()
	for _, q := range []Query{
		Type("blog"),
		Slug("Hello World"),
		Status(0),
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
	_post, err := PostFromBytes(metadata)
	if err != nil {
		t.Fatal(err)
	}

	toBe(t, "post.KV.title", string(post.KV["title"]), "Hello World!")
	toBe(t, "post.KV.date", string(post.KV["date"]), "2020-01-01")
	toBe(t, "_post.ID", _post.ID.String(), post.ID.String())
	toBe(t, "_post.Type", _post.Type, post.Type)
	toBe(t, "_post.ACL", _post.Status, post.Status)
	toBe(t, "_post.Owner", _post.Owner, post.Owner)
	toBe(t, "_post.Crtime", _post.Crtime, post.Crtime)
	toBe(t, "_post.Tags", strings.Join(_post.Tags, ","), strings.Join(post.Tags, ","))

	t.Log(_post)
}

func toBe(t *testing.T, name string, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("%s(%v) should equal to %v", name, a, b)
	}
}
