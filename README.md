# PostDB

[![GoDoc](https://godoc.org/github.com/postui/postdb?status.svg)](https://godoc.org/github.com/postui/postdb)
[![GoReport](https://goreportcard.com/badge/github.com/postui/postdb)](https://goreportcard.com/report/github.com/postui/postdb)
[![MIT](https://img.shields.io/badge/license-MIT-green)](./LICENSE)

A database to store posts in [Go](https://golang.org) with [BoltDB](https://github.com/etcd-io/bbolt), noSQL.

## Installation
```bash
go get github.com/postui/postdb
```

## Usage

```go
// opening a database
db, err := postdb.Open("post.db", 0666)
if err != nil {
    return err
}
defer db.Close()

// get all posts in the database
db.List()

// get posts with query
db.List(q.Tags("foo", "bar"), q.Range(1, 100), q.Order(q.DESC), q.Select("title", "date", "content"))

// get post without kv
db.Get(q.ID(id))

// get post with specified kv
db.Get(q.ID(id), q.Select("title", "date"))

// get post with prefixed kv
db.Get(q.ID(id), q.Select("title_*")) // match key in title_en,title_zh...

// get post with full kv
db.Get(q.ID(id), q.Select("*"))

// add a new post
db.Put(q.Alias(alias), q.Status(1), q.Tags("foo", "bar"), q.KV{"foo": []byte("bar")})

// update the existing post
db.Update(q.ID(id), q.KV{"foo": []byte("cool")})

// move the existing post
db.MoveTo(q.ID(id), q.Anchor(id))

// delete the existing post kv
db.DeleteKV(q.ID(id), q.Select("foo"))

// delete the existing posts
db.Delete(q.ID(id))

// using namespace
ns := db.Namespace("name")
ns.List()
ns.Get(q.ID(id))
...

// backup the entire database
db.WriteTo(w)
```
