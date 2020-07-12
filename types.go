package postdb

import (
	"io"

	"github.com/postui/postdb/q"
)

var (
	valuesKey    = []byte("values")
	postmetaKey  = []byte("postmeta")
	postindexKey = []byte("postindex")
	postkvKey    = []byte("postkv")
	slugsKey     = []byte("slugs")
	typesKey     = []byte("types")
	ownersKey    = []byte("owners")
	tagsKey      = []byte("tags")
)

type Database interface {
	Begin(writable bool) (*Tx, error)
	GetValue(key string) ([]byte, error)
	PutValue(key string, value []byte) error
	GetPosts(qs ...q.Query) ([]q.Post, error)
	GetPost(qs ...q.Query) (*q.Post, error)
	AddPost(qs ...q.Query) (*q.Post, error)
	UpdatePost(qs ...q.Query) error
	RemovePost(qs ...q.Query) error
	WriteTo(w io.Writer) (int64, error)
}

type NSDatabase interface {
	Namespace(name string) (Database, error)
	Database
}
