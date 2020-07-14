package postdb

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/postui/postdb/q"
)

func TestDB(t *testing.T) {
	db, err := Open("test.db", 0666)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	posts, err := db.GetPosts()
	if err != nil {
		t.Fatal(err)
	}

	post, err := db.AddPost(
		q.Type("news"),
		q.Slug(fmt.Sprintf("hello-world-%d", len(posts))),
		q.Status(2),
		q.Owner("admin"),
		q.Tags("hello", "world", "世界", "你好"),
		q.KV{
			"title": []byte(fmt.Sprintf("Hello World #%d", len(posts))),
			"date":  []byte(time.Now().Format(http.TimeFormat)),
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("post(%s,%s) added", post.ID, post.Slug)

	posts, err = db.GetPosts(q.Tags("世界"), q.Keys("title"), q.Order(q.ASC), q.Range("bs6pmhh8d3b520r8c05g", 6))
	if err != nil {
		t.Fatal(err)
	}
	t.Log("posts by tag(世界): ", len(posts))
}
