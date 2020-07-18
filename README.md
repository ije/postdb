# POSTDB

[![GoDoc](https://godoc.org/github.com/postui/postdb?status.svg)](https://godoc.org/github.com/postui/postdb)
[![GoReport](https://goreportcard.com/badge/github.com/postui/postdb)](https://goreportcard.com/report/github.com/postui/postdb)
[![MIT](https://img.shields.io/badge/license-MIT-green)](./LICENSE)

A database to store your posts in [golang](https://golang.org) with [bbolt](https://github.com/etcd-io/bbolt), noSQL.


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

// get the value for a key
db.GetValue("k")

// put the value for a key
db.PutValue("k", []byte("v"))

// get all posts
db.GetPosts()

// get posts with query
db.GetPosts(q.Type("type"), q.Tags("tag"), q.Limit(100), q.Order(q.DESC), q.Keys("title", "thumb"))

// get post by id without kv
db.GetPost(q.ID("id"))

// get post by id with specified kv
db.GetPost(q.ID("id"), q.Keys("title", "content"))

// get post by id with full kv
db.GetPost(q.ID("id"), q.Keys("*")))

// add a new post
db.AddPost(q.Type("type"), q.Slug("slug"), q.Tags("tag1", "tag2"), q.KV{"k": []byte("v1")})

// update the existing post
db.UpdatePost(q.ID("id"), q.KV{"k": []byte("v2")})

// remove the existing post permanently
db.RemovePost(q.ID("id"))

// backup the entire database
db.WriteTo(w)
```

with namespace:
```go
// opening a ns database
db, err := postdb.New("path")
if err != nil {
    return err
}
defer db.Close()
 
// use default namespace "public"
db.GetPosts()
...

// opening the namespace database
ns, err := db.Namespace("name")
if err != nil {
    return err
}

// use the namespace database
ns.GetPosts()
...
```

as server of C/S:

```go
// Opening a ns database
db, err := postdb.New("path")
if err != nil {
    return err
}
defer db.Close()

// start the server
server := &postdb.Server{
    DB:     db,
    Port:   9000,
    Secret: "PASS",
}
server.Serve()
```

as client of C/S:

```go
// connect to server
db, err := postdb.Connect(postdb.ConnConfig{
    Host:   "localhost"
    Port:   9000,
    Secret: "PASS",
})
if err != nil {
    return err
}

// use the database
db.GetPosts()
...
```


#   

Copyright (c) 2020-present, [postUI Inc.](https://postui.com)
