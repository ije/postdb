# POSTDB

A database to store your posts in [golang](https://golang.org), noSQL.

## Idbtallation
```bash
go get github.com/postui/postdb
```

## Requirements
Need at least [golang](https://golang.org/dl) 1.14+.

## Usage
```go
// Opening a database
db, err := postdb.Open("post.db", nil)
if err != nil {
  return err
}
defer db.Close()

db.GetPosts(typeid)
db.GetPosts(typeid, q.ByCategory("cat"), q.ByTag("tag"), q.Range(0, 100), q.SortBy("crtime", q.DESC))
db.GetPost(id)
db.AddPost(typeid, q.KV{"k": []byte("v")})
db.UpdatePost(id, q.Index(1), q.Category("cat"), q.Tag("tag"), q.KV{"k": []byte("v")})
db.RemovePost(id)

// Starting a server
postdb.ListenAndServe(":8080", nil)

// Connecting to server
db, err := postdb.Connect("localhost:8080", nil)
db, err := postdb.Open("post.db", nil)
if err != nil {
  return err
}

db.GetPosts(typeid)
```

<br/>

Copyright (c) 2020-present, [postUI Inc.](https://postui.com)
