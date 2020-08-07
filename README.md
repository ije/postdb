# POSTDB

[![GoDoc](https://godoc.org/github.com/postui/postdb?status.svg)](https://godoc.org/github.com/postui/postdb)
[![GoReport](https://goreportcard.com/badge/github.com/postui/postdb)](https://goreportcard.com/report/github.com/postui/postdb)
[![MIT](https://img.shields.io/badge/license-MIT-green)](./LICENSE)

A database to store your posts in [Golang](https://golang.org) with [BoltDB](https://github.com/etcd-io/bbolt), noSQL.


## Installing
```bash
go get github.com/postui/postdb
``` 


## Usage

as an embedded database:

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
db.List(q.Tags("tag"), q.Limit(100), q.Order(q.DESC), q.K("title", "thumb"))

// get post without kv
db.Get(q.ID("id"))

// get post with specified kv
db.Get(q.ID("id"), q.K("title", "content"))

// get post with prefixed kv
db.Get(q.ID("id"), q.K("title_*")) // match key in title_en,title_zh,title_de...

// get post with full kv
db.Get(q.ID("id"), q.K("*")))

// add a new post
db.Put(q.Alias("alias"), q.Tags("tag1", "tag2"), q.KV{"k": []byte("v")})

// update the existing post
db.Update(q.ID("id"), q.KV{"k2": []byte("v2")})

// move the existing post
db.MoveTo(q.ID("id"), q.Anchor("id"))

// delete the existing post kv
db.DeleteKV(q.ID("id"), q.K("k2"))

// delete the existing posts permanently
db.Delete(q.ID("id"))

// backup the entire database
db.WriteTo(w)
```

using namespace:

```go
// opening a database
db, err := postdb.Open("post.db", 0666)
if err != nil {
    return err
}
defer db.Close()

// creating a namespace
ns, err := db.Namespace("name")
if err != nil {
    return err
}

// use the db with namespace
ns.List()
...
```

as server of C/S:

```go
// create a server
s := &postdb.Server{
    DBPath: "path",
    Port:   9000,
}

// create a new database user
s.CreateUser("USERNAME", "PASS")

// start the server
s.Serve()
```

as client of C/S:

```go
// connect to server
client, err := postdb.Connect(postdb.ConnConfig{
    Host:     "localhost"
    Port:     9000,
    User:     "USERNAME",
    Password: "PASS",
})
if err != nil {
    return err
}

// opening a client database
db, err := client.DB("name")
if err != nil {
    return err
}

// use the client database
db.List()
...
```


#   

Copyright (c) 2020-present, [postUI Inc.](https://postui.com)
