package postdb

import (
	"bytes"
	"io"
	"strings"

	"github.com/postui/postdb/q"
	bolt "go.etcd.io/bbolt"
)

// Tx represents a transaction on the database.
type Tx struct {
	t *bolt.Tx
}

// List returns some posts
func (tx *Tx) List(qs ...q.Query) (posts []q.Post) {
	var cur *bolt.Cursor
	var prefixs [][]byte
	var filter func(*q.Post) bool
	var n uint32

	var res q.Resolver
	for _, q := range qs {
		res.Apply(q)
	}

	if res.BadID || res.BadAfter {
		return
	}

	queryTags := len(res.Tags) > 0
	queryOwner := len(res.Owner) > 0
	queryAfter := len(res.After) == 12
	orderASC := res.Order >= q.ASC

	metaBucket := tx.t.Bucket(postmetaKey)
	indexBucket := tx.t.Bucket(postindexKey)
	kvBucket := tx.t.Bucket(postkvKey)

	if len(res.ID) == 12 {
		v := metaBucket.Get(res.ID)
		if v != nil {
			post, err := q.PostFromBytes(v)
			if err == nil && (!queryOwner || post.Owner == res.Owner) {
				posts = make([]q.Post, 1)
				posts[0] = *post
			}
		}
	} else if len(res.Alias) > 0 {
		aliasIndexBucket := indexBucket.Bucket(postaliasKey)
		id := aliasIndexBucket.Get([]byte(res.Alias))
		if len(id) == 12 {
			v := metaBucket.Get(id)
			if v != nil {
				post, err := q.PostFromBytes(v)
				if err == nil && (!queryOwner || post.Owner == res.Owner) {
					posts = make([]q.Post, 1)
					posts[0] = *post
				}
			}
		}
	} else if queryTags {
		cur = indexBucket.Bucket(posttagKey).Cursor()
		prefixs = make([][]byte, len(res.Tags))
		var i int
		for tag := range res.Tags {
			p := make([]byte, len(tag)+1)
			copy(p, []byte(tag))
			prefixs[i] = p
			i++
		}
		if queryOwner {
			filter = func(post *q.Post) bool {
				return post.Owner == res.Owner
			}
		}
	} else if queryOwner {
		prefix := make([]byte, len(res.Owner)+1)
		copy(prefix, []byte(res.Owner))
		cur = indexBucket.Bucket(postownerKey).Cursor()
		prefixs = [][]byte{prefix}
	} else {
		c := metaBucket.Cursor()
		var k, v []byte
		if queryAfter {
			k, v = c.Seek(res.After)
		} else {
			if orderASC {
				k, v = c.First()
			} else {
				k, v = c.Last()
			}
		}
		for {
			if k == nil {
				break
			}
			post, err := q.PostFromBytes(v)
			if err == nil {
				posts = append(posts, *post)
				n++
				if res.Limit > 0 && n >= res.Limit {
					break
				}
			}
			if orderASC {
				k, v = c.Next()
			} else {
				k, v = c.Prev()
			}
		}
	}

	if cur != nil {
		pl := len(prefixs)
		a := make([][]*q.Post, pl)
		for i, prefix := range prefixs {
			k, _ := cur.Seek(prefix)
			if k == nil || !bytes.HasPrefix(k, prefix) {
				break
			}

			if orderASC && queryAfter {
				for {
					if bytes.Compare(k[len(k)-12:], res.After) <= 0 {
						k, _ = cur.Next()
						if k == nil || !bytes.HasPrefix(k, prefix) {
							break
						}
					} else {
						break
					}
				}
			}
			if !orderASC {
				for {
					k, _ = cur.Next()
					if k == nil {
						k, _ = cur.Last()
						break
					} else if !bytes.HasPrefix(k, prefix) || (queryAfter && bytes.Compare(k[len(k)-12:], res.After) >= 0) {
						k, _ = cur.Prev()
						break
					}
				}
			}

			for {
				if k == nil || !bytes.HasPrefix(k, prefix) {
					break
				}
				id := k[len(k)-12:]
				if i == 0 {
					data := metaBucket.Get(id)
					if data != nil {
						post, err := q.PostFromBytes(data)
						if err == nil && bytes.Compare(post.ID.Bytes(), id) == 0 && (filter == nil || filter(post) == true) && (!res.HasStatus || post.Status == res.Status) {
							if pl == 1 {
								posts = append(posts, *post)
							} else {
								a[0] = append(a[0], post)
							}

							n++
							if res.Limit > 0 && n >= res.Limit {
								break
							}
						}
					}
				} else if past := a[i-1]; len(past) > 0 {
					for _, p := range past {
						if bytes.Equal(p.ID.Bytes(), id) {
							if i == pl-1 {
								posts = append(posts, *p)
							} else {
								a[i] = append(a[i], p)
							}
						}
					}
				}
				if orderASC {
					k, _ = cur.Next()
				} else {
					k, _ = cur.Prev()
				}
			}
		}
	}

	if len(res.KVKeys) > 0 {
		for _, post := range posts {
			postkvBucket := kvBucket.Bucket(post.ID.Bytes())
			if postkvBucket != nil {
				if res.KVWildcard {
					c := postkvBucket.Cursor()
					for k, v := c.First(); k != nil; k, v = c.Next() {
						post.KV[string(k)] = v
					}
				} else {
					for key := range res.KVKeys {
						v := postkvBucket.Get([]byte(key))
						if v != nil {
							post.KV[key] = v
						}
					}
				}
			}
		}
	}

	return
}

// Get returns the post
func (tx *Tx) Get(qs ...q.Query) (*q.Post, error) {
	var res q.Resolver
	for _, q := range qs {
		res.Apply(q)
	}

	var metaData []byte
	if len(res.ID) == 12 {
		metaData = tx.t.Bucket(postmetaKey).Get(res.ID)
	} else if len(res.Alias) > 0 {
		aliasIndexBucket := tx.t.Bucket(postindexKey).Bucket(postaliasKey)
		id := aliasIndexBucket.Get([]byte(res.Alias))
		if len(id) == 12 {
			metaData = tx.t.Bucket(postmetaKey).Get(id)
		}
	}
	if metaData == nil {
		return nil, ErrNotFound
	}

	post, err := q.PostFromBytes(metaData)
	if err != nil {
		return nil, err
	}

	if len(res.KVKeys) > 0 {
		postkvBucket := tx.t.Bucket(postkvKey).Bucket(post.ID.Bytes())
		if postkvBucket != nil {
			if res.KVWildcard {
				c := postkvBucket.Cursor()
				for k, v := c.First(); k != nil; k, v = c.Next() {
					post.KV[string(k)] = v
				}
			} else {
				for key := range res.KVKeys {
					v := postkvBucket.Get([]byte(key))
					if v != nil {
						post.KV[key] = v
					}
				}
			}
		}
	}

	return post, nil
}

// Put puts a new post
func (tx *Tx) Put(qs ...q.Query) (*q.Post, error) {
	metaBucket := tx.t.Bucket(postmetaKey)
	indexBucket := tx.t.Bucket(postindexKey)
	kvBucket := tx.t.Bucket(postkvKey)

	post := q.NewPost()
	for _, q := range qs {
		post.ApplyQuery(q)
	}

	err := metaBucket.Put(post.ID.Bytes(), post.MetaData())
	if err != nil {
		return nil, err
	}

	if len(post.Alias) > 0 {
		aliasIndexBucket := indexBucket.Bucket(postaliasKey)
		if aliasIndexBucket.Get([]byte(post.Alias)) != nil {
			return nil, ErrDuplicateAlias
		}
		err = aliasIndexBucket.Put([]byte(post.Alias), post.ID.Bytes())
		if err != nil {
			return nil, err
		}
	}

	if len(post.Owner) > 0 {
		ownerIndexBucket := indexBucket.Bucket(postownerKey)
		keypath := [][]byte{[]byte(post.Owner), post.ID.Bytes()}
		err = ownerIndexBucket.Put(bytes.Join(keypath, []byte{0}), []byte{1})
		if err != nil {
			return nil, err
		}
	}

	if len(post.Tags) > 0 {
		tagIndexBucket := indexBucket.Bucket(posttagKey)
		for _, tag := range post.Tags {
			keypath := [][]byte{[]byte(tag), post.ID.Bytes()}
			err = tagIndexBucket.Put(bytes.Join(keypath, []byte{0}), []byte{1})
			if err != nil {
				return nil, err
			}
		}
	}

	postkvBucket, err := kvBucket.CreateBucketIfNotExists(post.ID.Bytes())
	if err != nil {
		return nil, err
	}

	if len(post.KV) > 0 {
		for k, v := range post.KV {
			err = postkvBucket.Put([]byte(k), v)
			if err != nil {
				return nil, err
			}
		}
	}

	return post, nil
}

// Update updates the post
func (tx *Tx) Update(qs ...q.Query) (*q.Post, error) {
	metaBucket := tx.t.Bucket(postmetaKey)
	indexBucket := tx.t.Bucket(postindexKey)
	kvBucket := tx.t.Bucket(postkvKey)

	var res q.Resolver
	for _, q := range qs {
		res.Apply(q)
	}

	var metaData []byte
	if len(res.ID) == 12 {
		metaData = metaBucket.Get(res.ID)
	}
	if metaData == nil {
		return nil, ErrNotFound
	}

	post, err := q.PostFromBytes(metaData)
	if err != nil {
		return nil, err
	}

	copy := post.Clone()
	for _, q := range qs {
		copy.ApplyQuery(q)
	}

	var shouldUpdateMeta bool

	// update alias index
	if copy.Alias != post.Alias {
		aliasIndexBucket := indexBucket.Bucket(postaliasKey)
		if aliasIndexBucket.Get([]byte(copy.Alias)) != nil {
			return nil, ErrDuplicateAlias
		}
		if len(post.Alias) > 0 {
			err = aliasIndexBucket.Delete([]byte(post.Alias))
			if err != nil {
				return nil, err
			}
		}
		if len(copy.Alias) > 0 {
			err = aliasIndexBucket.Put([]byte(copy.Alias), copy.ID.Bytes())
			if err != nil {
				return nil, err
			}
		}
		if !shouldUpdateMeta {
			shouldUpdateMeta = true
		}
	}

	// update owner index
	if copy.Owner != post.Owner {
		ownerIndexBucket := indexBucket.Bucket(postownerKey)
		if len(post.Owner) > 0 {
			keypath := [][]byte{[]byte(post.Owner), post.ID.Bytes()}
			err = ownerIndexBucket.Delete(bytes.Join(keypath, []byte{0}))
			if err != nil {
				return nil, err
			}
		}
		if len(copy.Owner) > 0 {
			keypath := [][]byte{[]byte(copy.Owner), copy.ID.Bytes()}
			err = ownerIndexBucket.Put(bytes.Join(keypath, []byte{0}), []byte{1})
			if err != nil {
				return nil, err
			}
		}
		if !shouldUpdateMeta {
			shouldUpdateMeta = true
		}
	}

	// update tags index
	if strings.Join(copy.Tags, "") != strings.Join(post.Tags, "") {
		tagIndexBucket := indexBucket.Bucket(posttagKey)
		if len(post.Tags) > 0 {
			for _, tag := range post.Tags {
				keypath := [][]byte{[]byte(tag), post.ID.Bytes()}
				err = tagIndexBucket.Delete(bytes.Join(keypath, []byte{0}))
				if err != nil {
					return nil, err
				}
			}
		}
		if len(copy.Tags) > 0 {
			for _, tag := range copy.Tags {
				keypath := [][]byte{[]byte(tag), copy.ID.Bytes()}
				err = tagIndexBucket.Put(bytes.Join(keypath, []byte{0}), []byte{1})
				if err != nil {
					return nil, err
				}
			}
		}
		if !shouldUpdateMeta {
			shouldUpdateMeta = true
		}
	}

	if shouldUpdateMeta {
		err = metaBucket.Put(copy.ID.Bytes(), copy.MetaData())
		if err != nil {
			return nil, err
		}
	}

	if len(copy.KV) > 0 {
		postkvBucket := kvBucket.Bucket(copy.ID.Bytes())
		for k, v := range copy.KV {
			err = postkvBucket.Put([]byte(k), v)
			if err != nil {
				return nil, err
			}
		}
	}

	return copy, nil
}

// DeleteKV deletes the post kv
func (tx *Tx) DeleteKV(qs ...q.Query) (err error) {
	post, err := tx.Get(qs...)
	if err != nil {
		return err
	}

	var res q.Resolver
	for _, q := range qs {
		res.Apply(q)
	}
	if len(res.KVKeys) > 0 {
		kvBucket := tx.t.Bucket(postkvKey)
		postkvBucket := kvBucket.Bucket(post.ID.Bytes())
		if postkvBucket != nil {
			for key := range res.KVKeys {
				err = postkvBucket.Delete([]byte(key))
				if err != nil {
					return
				}
			}
		}
	}

	return
}

// Delete deletes the post
func (tx *Tx) Delete(qs ...q.Query) (n int, err error) {
	if len(qs) == 0 {
		return 0, nil
	}

	metaBucket := tx.t.Bucket(postmetaKey)
	indexBucket := tx.t.Bucket(postindexKey)
	kvBucket := tx.t.Bucket(postkvKey)
	posts := tx.List(qs...)

	for _, post := range posts {
		err = metaBucket.Delete(post.ID.Bytes())
		if err != nil {
			return
		}

		if len(post.Alias) > 0 {
			aliasIndexBucket := indexBucket.Bucket(postaliasKey)
			err = aliasIndexBucket.Delete([]byte(post.Alias))
			if err != nil {
				return
			}
		}

		if len(post.Owner) > 0 {
			ownerIndexBucket := indexBucket.Bucket(postownerKey)
			keypath := [][]byte{[]byte(post.Owner), post.ID.Bytes()}
			err = ownerIndexBucket.Delete(bytes.Join(keypath, []byte{0}))
			if err != nil {
				return
			}
		}

		if len(post.Tags) > 0 {
			tagIndexBucket := indexBucket.Bucket(posttagKey)
			for _, tag := range post.Tags {
				keypath := [][]byte{[]byte(tag), post.ID.Bytes()}
				err = tagIndexBucket.Delete(bytes.Join(keypath, []byte{0}))
				if err != nil {
					return
				}
			}
		}

		err = kvBucket.DeleteBucket(post.ID.Bytes())
		if err != nil {
			return
		}
	}

	n = len(posts)
	return
}

// Rollback closes the transaction and ignores all previous updates. Read-only
// transactions must be rolled back and not committed.
func (tx *Tx) Rollback() error {
	return tx.t.Rollback()
}

// Commit writes all changes to disk and updates the meta page.
// Returns an error if a disk write error occurs, or if Commit is
// called on a read-only transaction.
func (tx *Tx) Commit() error {
	return tx.t.Commit()
}

// WriteTo writes the entire database to a writer.
// If err == nil then exactly tx.Size() bytes will be written into the writer.
func (tx *Tx) WriteTo(w io.Writer) (int64, error) {
	return tx.t.WriteTo(w)
}
