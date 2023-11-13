package postdb

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/ije/postdb/q"
)

func TestNS(t *testing.T) {
	db, err := Open("ns_test.db", 0666, false)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// flush
	_, err = db.Delete(q.Owner(42))
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.NS("a").Delete(q.Owner(42))
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.NS("b").Delete(q.Owner(42))
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 5; i++ {
		_, err := db.Put(
			q.Alias(fmt.Sprintf("hello-world-%d", i+1)),
			q.Status(1),
			q.Owner(42),
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
		_, err := db.NS("a").Put(
			q.Alias(fmt.Sprintf("hello-world-%d", i+1)),
			q.Status(1),
			q.Owner(42),
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

	_, err = db.NS("b").Put(
		q.Alias("hello-world-zh"),
		q.Status(1),
		q.Owner(42),
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

	posts, err = db.NS("a").List()
	if err != nil {
		t.Fatal(err)
	}
	toBe(t, "nsa's posts length", len(posts), 10)

	posts, err = db.NS("b").List()
	if err != nil {
		t.Fatal(err)
	}
	toBe(t, "nsb's posts length", len(posts), 1)
}
