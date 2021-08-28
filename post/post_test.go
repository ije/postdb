package post

import (
	"strings"
	"testing"
)

func TestPostMeta(t *testing.T) {
	post := New()
	post.Alias = "hello-world"
	post.Owner = "admin"
	post.Tags = []string{"hello", "world"}

	metadata := post.MetaData()
	_post, err := FromBytes(metadata)
	if err != nil {
		t.Fatal(err)
	}

	toBe(t, "_post.PKey", _post.PKey, post.PKey)
	toBe(t, "_post.ID", _post.ID, post.ID)
	toBe(t, "_post.Alias", _post.Alias, post.Alias)
	toBe(t, "_post.ACL", _post.Status, post.Status)
	toBe(t, "_post.Owner", _post.Owner, post.Owner)
	toBe(t, "_post.Crtime", _post.Crtime, post.Crtime)
	toBe(t, "_post.Tags", strings.Join(_post.Tags, ","), strings.Join(post.Tags, ","))
}

func toBe(t *testing.T, name string, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("the %s should equal to %v, but %v", name, b, a)
	}
}
