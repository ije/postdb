package postdb

import (
	"io"

	"github.com/postui/postdb/q"
)

type ClientDB struct {
	db     string
	client *Client
}

// Begin starts a new transaction.
func (db *ClientDB) Begin(writable bool) (*ClientTx, error) {
	return &ClientTx{db: db}, nil
}

// List returns some posts
func (db *ClientDB) List(qs ...q.Query) ([]q.Post, error) {
	tx, err := db.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	return tx.List(qs...), nil
}

// Get returns the post
func (db *ClientDB) Get(qs ...q.Query) (*q.Post, error) {
	tx, err := db.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	return tx.Get(qs...)
}

// Put puts a new post
func (db *ClientDB) Put(qs ...q.Query) (*q.Post, error) {
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
func (db *ClientDB) Update(qs ...q.Query) error {
	tx, err := db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = tx.Update(qs...)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// DeleteKV deletes the post kv
func (db *ClientDB) DeleteKV(qs ...q.Query) (int, error) {
	tx, err := db.Begin(true)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	n, err := tx.DeleteKV(qs...)
	if err != nil {
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return n, nil
}

// Delete deletes the post
func (db *ClientDB) Delete(qs ...q.Query) (int, error) {
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
func (db *ClientDB) WriteTo(w io.Writer) (int64, error) {
	tx, err := db.Begin(false)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	return tx.WriteTo(w)
}

// Close releases all database resources.
func (db *ClientDB) Close() error {
	return nil
}
