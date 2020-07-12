# POSTDB

A database to store your posts in [golang](https://golang.org) with [bbolt](https://github.com/etcd-io/bbolt), noSQL.

[![GoDoc](https://godoc.org/github.com/postui/postdb?status.svg)](https://godoc.org/github.com/postui/postdb)

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
db.GetPosts(q.Type("type"), q.Tags("tag"), q.Range("", 100), q.Order(q.DESC))

// get post by id without kv
db.GetPost(q.ID("id"))

// get post by id with kv
db.GetPost(q.ID("id"), q.Keys{"title", "thumb", "content"))

// add a new post
db.AddPost(q.Type("type"), q.Slug("slug"), q.Tags("tag"), q.KV{"k": []byte("v")})

// update the existing post
db.UpdatePost(q.ID("id"),  q.Tags("tagA", "tagB"), q.KV{"k": []byte("v")})

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

// create namespace database
ns := db.Namespace("name")

// use namespace database
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

// use database
db.GetPosts()
...
```

<br/>

__

Copyright (c) 2020-present, [postUI Inc.](https://postui.com)
