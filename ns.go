package postdb

import (
	"github.com/ije/postdb/internal/post"
	"github.com/ije/postdb/q"
)

// A NS for the DB
type NS struct {
	err  error
	name string
	db   *DB
}

// Begin starts a new transaction.
func (ns *NS) Begin(writable bool) (*Tx, error) {
	if ns.err != nil {
		return nil, ns.err
	}

	tx, err := ns.db.bolt.Begin(writable)
	if err != nil {
		return nil, err
	}

	return &Tx{[]byte(ns.name), tx}, nil
}

// List returns some posts
func (ns *NS) List(qs ...q.Query) ([]post.Post, error) {
	tx, err := ns.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	return tx.List(qs...), nil
}

// Get returns the post
func (ns *NS) Get(qs ...q.Query) (*post.Post, error) {
	tx, err := ns.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	return tx.Get(qs...)
}

// Put puts a new post
func (ns *NS) Put(qs ...q.Query) (*post.Post, error) {
	tx, err := ns.Begin(true)
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
func (ns *NS) Update(qs ...q.Query) error {
	tx, err := ns.Begin(true)
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
func (ns *NS) DeleteKV(qs ...q.Query) error {
	tx, err := ns.Begin(true)
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

// MoveTo moves the post
func (ns *NS) MoveTo(qs ...q.Query) error {
	tx, err := ns.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = tx.MoveTo(qs...)
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
func (ns *NS) Delete(qs ...q.Query) (int, error) {
	tx, err := ns.Begin(true)
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
