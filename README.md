# POSTDB

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

// get typed posts
db.GetPosts("type")
// get typed posts with query
db.GetPosts("type", q.Status(q.DRAFT, q.NORMAL), q.Tag("tag"), q.Range(0, 100), q.SortBy("crtime", q.DESC))
// get post by id
db.GetPost("id")
// add a new post
db.AddPost("type", q.Tag("tag"), q.KV{"k": []byte("v")})
// update the existing post
db.UpdatePost("id", q.Tag("tagA", "tagB"), q.KV{"k": []byte("v")})
// remove the existing post
db.RemovePost("id")
// backup the database
db.WriteTo(w)

// get the value for a key
db.GetValue("k")
// put the value for a key
db.PutValue("k", []byte("v"))
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
db.GetPosts("type")
...

// create namespace database
nsdb := db.Namespace("name")
// use namespace database
nsdb.GetPosts("type")
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

// start a server
postdb.ListenAndServe(db, &postdb.ServerConfig{
  Port: 9000,
  Secret: "PASS",
})
```

as client of C/S:

```go
// connect to server
db, err := postdb.Connect(&postdb.ConnConfig{
    Host: "localhost"
    Port: 9000,
    Secret: "PASS",
})
if err != nil {
  return err
}

// use database
db.GetPosts("type")
...
```


as graphql http handler:

```go
// opening a ns database
db, err := postdb.New("path")
if err != nil {
  return err
}
defer db.Close()

graphql := postdb.NewGraphql(db)
// register graphql http handler
http.Handle("/graphql", graphql)
// with simple basic auth
http.Handle("/graphql", httpauth.SimpleBasicAuth("username", "PASS")(graphql))
http.ListenAndServe(":8080", nil)
```
<br/>

Copyright (c) 2020-present, [postUI Inc.](https://postui.com)
