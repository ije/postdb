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
		for _, key := range [][]byte{
			postmetaKey,
			postindexKey,
			postkvKey,
		} {
			_, err := tx.CreateBucketIfNotExists(key)
			if err != nil {
				return err
			}
		}
		indexBucket := tx.Bucket(postindexKey)
		for _, key := range [][]byte{
			postuidKey,
			postownerKey,
			posttagKey,
		} {
			_, err := indexBucket.CreateBucketIfNotExists(key)
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

	return &Tx{tx}, nil
}

// List returns some posts
func (db *DB) List(qs ...q.Query) ([]q.Post, error) {
	tx, err := db.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	return tx.List(qs...), nil
}

// Get returns the post
func (db *DB) Get(qs ...q.Query) (*q.Post, error) {
	tx, err := db.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	return tx.Get(qs...)
}

// Put puts a new post
func (db *DB) Put(qs ...q.Query) (*q.Post, error) {
	tx, err := db.Begin(true)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	post, err := tx.Put(qs...)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return post, nil
}

// Update updates the post
func (db *DB) Update(qs ...q.Query) (*q.Post, error) {
	tx, err := db.Begin(true)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	post, err := tx.Update(qs...)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return post, nil
}

// DeleteKV deletes the post kv
func (db *DB) DeleteKV(qs ...q.Query) error {
	tx, err := db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = tx.DeleteKV(qs...)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// Delete deletes the post
func (db *DB) Delete(qs ...q.Query) (int, error) {
	tx, err := db.Begin(true)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	n, err := tx.Delete(qs...)
	if err != nil {
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return n, nil
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
