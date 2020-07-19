package postdb

import (
	"io"

	"github.com/postui/postdb/q"
)

// ClientTx represents a transaction on the database.
type ClientTx struct {
	db *ClientDB
}

// List returns some posts
func (tx *ClientTx) List(qs ...q.Query) (posts []q.Post) {
	return nil
}

// Get returns the post
func (tx *ClientTx) Get(qs ...q.Query) (*q.Post, error) {
	return nil, nil
}

// Put puts a new post
func (tx *ClientTx) Put(qs ...q.Query) (*q.Post, error) {

	return nil, nil
}

// Update updates the post
func (tx *ClientTx) Update(qs ...q.Query) (*q.Post, error) {
	return nil, nil
}

// DeleteKV deletes the post kv
func (tx *ClientTx) DeleteKV(qs ...q.Query) (n int, err error) {
	return
}

// Delete deletes the post
func (tx *ClientTx) Delete(qs ...q.Query) (n int, err error) {
	return
}

// Rollback closes the transaction and ignores all previous updates. Read-only
// transactions must be rolled back and not committed.
func (tx *ClientTx) Rollback() error {
	return nil
}

// Commit writes all changes to disk and updates the meta page.
// Returns an error if a disk write error occurs, or if Commit is
// called on a read-only transaction.
func (tx *ClientTx) Commit() error {
	return nil
}

// WriteTo writes the entire database to a writer.
// If err == nil then exactly tx.Size() bytes will be written into the writer.
func (tx *ClientTx) WriteTo(w io.Writer) (int64, error) {
	return 0, nil
}
