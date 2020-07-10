package postdb

import (
	"io"
	"os"
	"time"

	"github.com/postui/postdb/q"
	bolt "go.etcd.io/bbolt"
)

// A DB to store posts
type DB struct {
	b *bolt.DB
}

// Open opens a database at the given path.
func Open(path string, mode os.FileMode) (db *DB, err error) {
	b, err := bolt.Open(path, mode, &bolt.Options{
		Timeout: 1 * time.Second,
	})
	if err != nil {
		return
	}

	err = b.Update(func(tx *bolt.Tx) error {
		bucketNames := [][]byte{valuesKey, postmetaKey, postkvKey}
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

// Begin starts a new transaction.
func (db *DB) Begin(writable bool) (*Tx, error) {
	tx, err := db.b.Begin(writable)
	if err != nil {
		return nil, err
	}
	return &Tx{tx: tx}, nil
}

// GetValue returns the value for a key in the database.
func (db *DB) GetValue(key string) ([]byte, error) {
	tx, err := db.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	return tx.GetValue(key)
}

// PutValue sets the value for a key in the database.
func (db *DB) PutValue(key string, value []byte) error {
	tx, err := db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = tx.PutValue(key, value)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (db *DB) GetPost(id string) (*Post, error) {
	tx, err := db.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	return tx.GetPost(id)
}

func (db *DB) GetPosts(postType string, qs ...[]q.Query) ([]Post, error) {
	tx, err := db.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	return tx.GetPosts(postType, qs...)
}

func (db *DB) AddPost(postType string, qs ...[]q.Query) error {
	tx, err := db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	return tx.AddPost(postType, qs...)
}

func (db *DB) UpdatePost(id string, qs ...[]q.Query) error {
	tx, err := db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	return tx.UpdatePost(id, qs...)
}

func (db *DB) RemovePost(id string) error {
	tx, err := db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	return tx.RemovePost(id)
}

// WriteTo writes the entire database to a writer.
func (db *DB) WriteTo(w io.Writer) (int64, error) {
	tx, err := db.Begin(false)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	return tx.WriteTo(w)
}

// Close releases all database resources.
func (db *DB) Close() error {
	return db.b.Close()
}
