package q

import (
	"strings"
	"testing"
)

func TestPostMeta(t *testing.T) {
	post := NewPost()
	for _, q := range []Query{
		Alias("Hello World"),
		Status(1),
		Owner("admin"),
		Tags("hello", "world"),
		KV{
			"title": []byte("Hello World!"),
			"date":  []byte("2020-01-01"),
		},
	} {
		q.Apply(post)
	}

	metadata := post.MetaData()
	_post, err := PostFromBytes(metadata)
	if err != nil {
		t.Fatal(err)
	}

	toBe(t, "post.KV.title", string(post.KV["title"]), "Hello World!")
	toBe(t, "post.KV.date", string(post.KV["date"]), "2020-01-01")
	toBe(t, "_post.PKey", _post.PKey, post.PKey)
	toBe(t, "_post.ID", _post.ID, post.ID)
	toBe(t, "_post.Alias", _post.Alias, post.Alias)
	toBe(t, "_post.ACL", _post.Status, post.Status)
	toBe(t, "_post.Owner", _post.Owner, post.Owner)
	toBe(t, "_post.Crtime", _post.Crtime, post.Crtime)
	toBe(t, "_post.Tags", strings.Join(_post.Tags, ","), strings.Join(post.Tags, ","))

	t.Log(_post)
}

func toBe(t *testing.T, name string, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("the %s should equal to %v, but %v", name, b, a)
	}
}
