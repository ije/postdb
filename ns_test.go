package postdb

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/postui/postdb/q"
)

func TestNS(t *testing.T) {
	db, err := Open("ns_test.db", 0666)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	nsa, err := db.Namespace("a")
	if err != nil {
		t.Fatal(err)
	}
	nsb, err := db.Namespace("b")
	if err != nil {
		t.Fatal(err)
	}

	// flush
	_, err = db.Delete(q.Owner("admin"))
	if err != nil {
		t.Fatal(err)
	}
	_, err = nsa.Delete(q.Owner("admin"))
	if err != nil {
		t.Fatal(err)
	}
	_, err = nsb.Delete(q.Owner("admin"))
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 5; i++ {
		_, err := db.Put(
			q.Alias(fmt.Sprintf("hello-world-%d", i+1)),
			q.Status(1),
			q.Owner("admin"),
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

	for i := 0; i < 10; i++ {
		_, err := nsa.Put(
			q.Alias(fmt.Sprintf("hello-world-%d", i+1)),
			q.Status(1),
			q.Owner("admin"),
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

	_, err = nsb.Put(
		q.Alias("hello-world-zh"),
		q.Status(1),
		q.Owner("admin"),
		q.Tags("你好", "世界"),
		q.KV{
			"title": []byte("你好世界！"),
			"date":  []byte(time.Now().Format(http.TimeFormat + " :)")),
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	posts, err := db.List()
	if err != nil {
		t.Fatal(err)
	}
	toBe(t, "db posts length", len(posts), 5)

	posts, err = nsa.List()
	if err != nil {
		t.Fatal(err)
	}
	toBe(t, "nsa's posts length", len(posts), 10)

	posts, err = nsb.List()
	if err != nil {
		t.Fatal(err)
	}
	toBe(t, "nsb's posts length", len(posts), 1)
}
