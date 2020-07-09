package postdb

import (
	"io"
	"os"
	"time"

	"github.com/postui/postdb/q"
	"github.com/postui/postdb/types"
	bolt "go.etcd.io/bbolt"
)

// A DB to store posts
type DB struct {
	b *bolt.DB
}

// Open opens a database at the given path.
func Open(path string, mode os.FileMode) (db *DB, err error) {
	b, err := bolt.Open(path, mode, &bolt.Options{
		Timeout: time.Second,
	})
	if err != nil {
		return
	}

	err = b.Update(func(tx *bolt.Tx) error {
		bucketNames := []string{"values", "postmeta", "postkv"}
		for _, name := range bucketNames {
			_, err := tx.CreateBucketIfNotExists([]byte(name))
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return
	}

	db = &DB{b}
	return
}

func (db *DB) GetValue(key string) ([]byte, error) {
	tx, err := db.b.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	return tx.Bucket([]byte("values")).Get([]byte(key)), nil
}

func (db *DB) PutValue(key string, value []byte) error {
	return db.b.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte("values")).Put([]byte(key), value)
	})
}

func (db *DB) GetPosts(postType string, qs ...[]q.Query) ([]types.Post, error) {
	return nil, nil
}

func (db *DB) GetPost(id string) (*types.Post, error) {
	return nil, nil
}

func (db *DB) AddPost(postType string, qs ...[]q.Query) (*types.Post, error) {
	return nil, nil
}

func (db *DB) UpdatePost(id string, qs ...[]q.Query) (*types.Post, error) {
	return nil, nil
}

func (db *DB) RemovePost(id string) error {
	return nil
}

func (db *DB) WriteTo(w io.Writer) (int64, error) {
	return 0, nil
}

func (db *DB) Close() error {
	return nil
}
