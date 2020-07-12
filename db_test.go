package postdb

import (
	"testing"

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

	t.Log("restore posts: ")
	for _, post := range posts {
		t.Log(post)
	}

	post, err := db.AddPost(
		q.Type("news"),
		q.Slug("hello-world"),
		q.Status(2),
		q.Owner("admin"),
		q.Tags("hello", "world", "世界"),
		q.KV{
			"title": []byte("Hello World!"),
			"date":  []byte("2020-01-01"),
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("post(%x) added", post.ID)

	posts, err = db.GetPosts(q.Tags("世界"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log("posts by tag(世界): ", len(posts))
}
