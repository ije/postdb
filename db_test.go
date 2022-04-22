package postdb

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/ije/postdb/q"
)

func TestDB(t *testing.T) {
	db, err := Open("test.db", 0666, false)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// flush
	_, err = db.Delete(q.Owner("tester"))
	if err != nil {
		t.Fatal(err)
	}

	posts, err := db.List(q.Owner("tester"))
	if err != nil {
		t.Fatal(err)
	}
	toBe(t, "posts length", len(posts), 0)

	for i := 0; i < 10; i++ {
		_, err := db.Put(
			q.Alias(fmt.Sprintf("hello-world-%d", i+1)),
			q.Status(1),
			q.Owner("tester"),
			q.Tags("hello", "world"),
			q.KV{
				"title": []byte(fmt.Sprintf("Hello World #%d", i+1)),
				"date":  []byte(time.Now().Format(http.TimeFormat)),
			},
		)
		if err != nil {
			t.Fatal(err)
		}
	}

	tp, err := db.Put(
		q.Alias("tmp"),
		q.Owner("tester"),
		q.KV{
			"title": []byte("Hello World!"),
			"date":  []byte(time.Now().Format(http.TimeFormat)),
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	posts, err = db.List()
	if err != nil {
		t.Fatal(err)
	}
	toBe(t, "posts length", len(posts), 11)

	_, err = db.Delete(q.ID(tp.ID))
	if err != nil {
		t.Fatal(err)
	}

	posts, err = db.List()
	if err != nil {
		t.Fatal(err)
	}
	toBe(t, "posts length", len(posts), 10)

	postZh, err := db.Put(
		q.Alias("hello-world-cn"),
		q.Status(1),
		q.Owner("abc"),
		q.Tags("hello", "world"),
		q.KV{
			"title": []byte("Hello World!"),
			"date":  []byte(time.Now().Format(http.TimeFormat)),
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	posts, err = db.List(q.Tags("world"))
	if err != nil {
		t.Fatal(err)
	}
	toBe(t, "posts length", len(posts), 11)

	err = db.Update(
		q.ID(postZh.ID),
		q.Alias("hello-world-zh"),
		q.Status(2),
		q.Owner("tester"),
		q.Tags("你好", "世界"),
		q.KV{
			"title":  []byte("你好世界！"),
			"date":   []byte(time.Now().Format(http.TimeFormat + " :)")),
			"newkey": []byte("v"),
		},
	)
	if err != nil {
		t.Fatal(postZh.ID, err)
	}

	_, err = db.Get(q.Alias("hello-world-zh"))
	if err != nil {
		t.Fatal(err)
	}

	postZh, err = db.Get(q.ID(postZh.ID), q.Select("*"))
	if err != nil {
		t.Fatal(err)
	}
	toBe(t, "postZh.Alias", postZh.Alias, "hello-world-zh")
	toBe(t, "postZh.Status", postZh.Status, uint8(2))
	toBe(t, "postZh.Owner", postZh.Owner, "tester")
	toBe(t, "postZh.Tags", strings.Join(postZh.Tags, " "), "你好 世界")
	toBe(t, "postZh.KV.title", string(postZh.KV["title"]), "你好世界！")
	toBe(t, "postZh.KV.date", strings.HasSuffix(string(postZh.KV["date"]), ":)"), true)
	toBe(t, "postZh.KV.newkey", string(postZh.KV["newkey"]), "v")

	posts, err = db.List()
	if err != nil {
		t.Fatal(err)
	}
	toBe(t, "posts length", len(posts), 11)

	posts, err = db.List(q.Tags("world"), q.Limit(5), q.Order(q.ASC), q.Select("*"))
	if err != nil {
		t.Fatal(err)
	}
	toBe(t, "posts length", len(posts), 5)
	for _, post := range posts {
		t.Logf(`%s/%s "%s" %s`, post.ID, post.Alias, string(post.KV["title"]), string(post.KV["date"]))
	}

	posts, err = db.List(q.Offset(10), q.Limit(5))
	if err != nil {
		t.Fatal(err)
	}
	toBe(t, "posts length", len(posts), 1)

	posts, err = db.List(q.Tags("世界"), q.Select("*"))
	if err != nil {
		t.Fatal(err)
	}
	toBe(t, "posts length", len(posts), 1)
	for _, post := range posts {
		t.Logf(`%s/%s "%s" %s`, post.ID, post.Alias, string(post.KV["title"]), string(post.KV["date"]))
	}
}

func toBe(t *testing.T, name string, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("the %s should equal to %v, but %v", name, b, a)
	}
}
