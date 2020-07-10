package postdb

import (
	"io"

	"github.com/postui/postdb/q"
	bolt "go.etcd.io/bbolt"
)

// Tx represents a transaction on the database.
type Tx struct {
	valuesBkt   *bolt.Bucket
	postmetaBkt *bolt.Bucket
	postkvBkt   *bolt.Bucket
	tx          *bolt.Tx
}

func (tx *Tx) GetValue(key string) ([]byte, error) {
	return tx.valuesBucket().Get([]byte(key)), nil
}

func (tx *Tx) PutValue(key string, value []byte) error {
	return tx.valuesBucket().Put([]byte(key), value)
}

func (tx *Tx) GetPosts(postType string, qs ...[]q.Query) ([]Post, error) {
	return nil, nil
}

func (tx *Tx) GetPost(id string) (*Post, error) {
	return nil, nil
}

func (tx *Tx) AddPost(postType string, qs ...[]q.Query) error {
	return nil
}

func (tx *Tx) UpdatePost(id string, qs ...[]q.Query) error {
	return nil
}

func (tx *Tx) RemovePost(id string) error {
	return nil
}

func (tx *Tx) valuesBucket() *bolt.Bucket {
	if tx.valuesBkt == nil {
		tx.valuesBkt = tx.tx.Bucket(valuesKey)
	}
	return tx.valuesBkt
}

func (tx *Tx) postmetaBucket() *bolt.Bucket {
	if tx.postmetaBkt == nil {
		tx.postmetaBkt = tx.tx.Bucket(postmetaKey)
	}
	return tx.postmetaBkt
}

func (tx *Tx) postkvBucket() *bolt.Bucket {
	if tx.postkvBkt == nil {
		tx.postkvBkt = tx.tx.Bucket(postkvKey)
	}
	return tx.postkvBkt
}

func (tx *Tx) clean() {
	if tx.valuesBkt != nil {
		tx.valuesBkt = nil
	}
	if tx.postmetaBkt != nil {
		tx.postmetaBkt = nil
	}
	if tx.postkvBkt != nil {
		tx.postkvBkt = nil
	}
}

// Rollback closes the transaction and ignores all previous updates. Read-only
// transactions must be rolled back and not committed.
func (tx *Tx) Rollback() error {
	tx.clean()
	return tx.tx.Rollback()
}

// Commit writes all changes to disk and updates the meta page.
// Returns an error if a disk write error occurs, or if Commit is
// called on a read-only transaction.
func (tx *Tx) Commit() error {
	tx.clean()
	return tx.tx.Commit()
}

// WriteTo writes the entire database to a writer.
// If err == nil then exactly tx.Size() bytes will be written into the writer.
func (tx *Tx) WriteTo(w io.Writer) (int64, error) {
	return tx.tx.WriteTo(w)
}
