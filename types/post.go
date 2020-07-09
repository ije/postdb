package types

import (
	"github.com/postui/postdb/q"
)

// A Post specifies a post of postdb
type Post struct {
	ID             string
	Type           string
	Tags           []string
	Status         PostStatus
	Crtime         int64
	Mtime          int64
	CurrentVersion uint32
	Content        q.KV
}
