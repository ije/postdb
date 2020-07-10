package postdb

import (
	"io"

	"github.com/postui/postdb/q"
)

var (
	valuesKey   = []byte("values")
	postmetaKey = []byte("postmeta")
	postkvKey   = []byte("postkv")
)

type Database interface {
	Begin(writable bool) (*Tx, error)
	GetValue(key string) ([]byte, error)
	PutValue(key string, value []byte) error
	GetPost(id string) (*Post, error)
	GetPosts(postType string, qs ...[]q.Query) ([]Post, error)
	AddPost(postType string, qs ...[]q.Query) error
	UpdatePost(id string, qs ...[]q.Query) error
	RemovePost(id string) error
	WriteTo(w io.Writer) (int64, error)
}

type NSDatabase interface {
	Namespace(name string) (Database, error)
	Database
}

// A Post specifies a post of postdb
type Post struct {
	ID             string
	Type           string
	Tags           []string
	Status         q.PostStatus
	Crtime         int64
	Mtime          int64
	CurrentVersion uint32
	Content        q.KV
}
