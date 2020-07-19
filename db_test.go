package postdb

import (
	"fmt"
	"net/http"
	"strings"
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

	// flush
	_, err = db.Delete(q.Type("test"))
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		_, err := db.Put(
			q.Type("test"),
			q.Slug(fmt.Sprintf("hello-world-%d", i+1)),
			q.Status(1),
			q.Owner("admin"),
			q.Tags("hello", "world", "世界"),
			q.KV{
				"title": []byte(fmt.Sprintf("Hello World #%d", i+1)),
				"date":  []byte(time.Now().Format(http.TimeFormat)),
			},
		)
		if err != nil {
			t.Fatal(err)
		}
	}

	postZh, err := db.Put(
		q.Type("test"),
		q.Slug("hello-world-cn"),
		q.Status(1),
		q.Owner("admin"),
		q.Tags("hello", "world", "世界"),
		q.KV{
			"title": []byte("Hello World!"),
			"date":  []byte(time.Now().Format(http.TimeFormat)),
			"k":     []byte("v1"),
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	posts, err := db.List(q.Tags("world"))
	if err != nil {
		t.Fatal(err)
	}
	toBe(t, "posts len", len(posts), 11)

	_, err = db.Update(
		postZh.ID,
		q.Slug("hello-world-zh"),
		q.Status(2),
		q.Owner("adminisitor"),
		q.Tags("你好", "世界"),
		q.KV{
			"title": []byte("你好世界！"),
			"date":  []byte(time.Now().Format(http.TimeFormat + " :)")),
			"k":     []byte("v2"),
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	postZh, err = db.Get(postZh.ID, q.Keys("*"))
	if err != nil {
		t.Fatal(err)
	}
	toBe(t, "postZh.Slug", postZh.Slug, "hello-world-zh")
	toBe(t, "postZh.Status", postZh.Status, uint8(2))
	toBe(t, "postZh.Owner", postZh.Owner, "adminisitor")
	toBe(t, "postZh.Tags", strings.Join(postZh.Tags, " "), "你好 世界")
	toBe(t, "postZh.KV.title", string(postZh.KV["title"]), "你好世界！")
	toBe(t, "postZh.KV.date", strings.HasSuffix(string(postZh.KV["date"]), ":)"), true)
	toBe(t, "postZh.KV.k", string(postZh.KV["k"]), "v2")

	posts, err = db.List(q.Tags("world"), q.Limit(5), q.Order(q.ASC), q.Keys("*"))
	if err != nil {
		t.Fatal(err)
	}
	toBe(t, "posts len", len(posts), 5)
	for i, post := range posts {
		t.Logf(`%d. %s/%s "%s" %s`, i+1, post.ID, post.Slug, string(post.KV.Get("title")), string(post.KV.Get("date")))
	}
}

func toBe(t *testing.T, name string, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("the %s should equal to %v, but %v", name, b, a)
	}
}
