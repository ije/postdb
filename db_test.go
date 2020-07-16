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
		q.Tags("hello", "world", "世界"),
		q.KV{
			"title": []byte(fmt.Sprintf("Hello World #%d", len(posts))),
			"date":  []byte(time.Now().Format(http.TimeFormat)),
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("post(%s/%s) added", post.ID, post.Slug)

	post, err = db.GetPost(q.ID(post.ID.String()), q.Keys("title"))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("restore post(%s/%s): %s", post.ID, post.Slug, string(post.KV.Get("title")))

	posts, err = db.GetPosts(q.Tags("世界"), q.Order(q.ASC), q.Limit(5), q.Keys("title"))
	if err != nil {
		t.Fatal(err)
	}
	for i, post := range posts {
		t.Logf(`%d. %s/%s "%s"`, i+1, post.ID, post.Slug, string(post.KV.Get("title")))
	}
}
