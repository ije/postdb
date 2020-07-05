# POSTDB

A database to store your posts in [golang](https://golang.org) with [bbolt](https://github.com/etcd-io/bbolt), noSQL.

## Installing
```bash
go get github.com/postui/postdb
```

## Requirements
Need [golang](https://golang.org/dl) 1.14+

## Usage
```go
// Opening a database
db, err := postdb.Open("post.db", nil)
if err != nil {
  return err
}
defer db.Close()

db.GetPosts(typeid)
db.GetPosts(typeid, q.ByCategory("cat"), q.Range(0, 100), q.SortBy("crtime", q.DESC))
db.GetPost(id)
db.AddPost(typeid, q.KV{"k": []byte("v")})
db.UpdatePost(id, q.Index(1), q.Tag("tag"), q.KV{"k": []byte("v")})
db.RemovePost(id)
```

as server:

```go
// Open a database
db, err := postdb.Open("post.db", nil)
if err != nil {
  return err
}
defer db.Close()

// Start a server
db.ListenAndServe(":9000", nil)
```

as client:

```go
// Connect to server
db, err := postdb.Connect("localhost:9000", nil)
if err != nil {
  return err
}

db.GetPosts(typeid)
```
<br/>

Copyright (c) 2020-present, [postUI Inc.](https://postui.com)
