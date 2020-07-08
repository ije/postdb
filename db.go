package postdb

import (
	"io"

	"github.com/postui/postdb/q"
	"github.com/postui/postdb/types"
)

// A DB to store posts
type DB interface {
	GetValue(key string) ([]byte, error)
	PutValue(key string, value []byte) error
	GetPosts(postType string, qs ...[]q.Query) ([]types.Post, error)
	GetPost(id string) (types.Post, error)
	AddPost(postType string, qs ...[]q.Query) (types.Post, error)
	UpdatePost(id string, qs ...[]q.Query) (types.Post, error)
	RemovePost(id string) error
	WriteTo(w io.Writer) (int64, error)
	Close() error
}

// A NSDB to store posts
type NSDB interface {
	Namespace(name string) (DB, error)
	DB
}
