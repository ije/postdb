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
	metaBucket := tx.t.Bucket(postmetaKey)
	indexBucket := tx.t.Bucket(postindexKey)
	kvBucket := tx.t.Bucket(postkvKey)
	uidIndexBucket := indexBucket.Bucket(postuidKey)

	var cur *bolt.Cursor
	var prefixs [][]byte
	var filter func(*q.Post) bool
	var n uint32

	var res q.Resolver
	for _, q := range qs {
		q.Resolve(&res)
	}

	queryOwner := len(res.Owner) > 0
	orderASC := res.Order >= q.ASC

	var afterPkey []byte
	if len(res.After) > 0 {
		afterPkey = uidIndexBucket.Get([]byte(res.After))
		if afterPkey == nil {
			return
		}
	}

	if len(res.ID) > 0 {
		pkey := uidIndexBucket.Get([]byte(res.ID))
		if pkey != nil {
			v := metaBucket.Get(pkey)
			if v != nil {
				post, err := q.PostFromBytes(v)
				if err == nil && bytes.Equal(post.PKey.Bytes(), pkey) && (!queryOwner || post.Owner == res.Owner) {
					posts = make([]q.Post, 1)
					posts[0] = *post
				}
			}
		}
	} else if len(res.IDs) > 0 {
		posts = make([]q.Post, len(res.IDs))
		for i, id := range res.IDs {
			pkey := uidIndexBucket.Get([]byte(id))
			if pkey != nil {
				v := metaBucket.Get(pkey)
				if v != nil {
					post, err := q.PostFromBytes(v)
					if err == nil && bytes.Equal(post.PKey.Bytes(), pkey) && (!queryOwner || post.Owner == res.Owner) {
						posts[i] = *post
						n++
						if res.Limit > 0 && n >= res.Limit {
							break
						}
					}
				}
			}
		}
		posts = posts[:n]
	} else if len(res.Tags) > 0 {
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
		if afterPkey != nil {
			k, v = c.Seek(afterPkey)
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
			if err == nil && bytes.Equal(post.PKey.Bytes(), k) && (!res.HasStatus || post.Status == res.Status) {
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
			ok := func(k []byte) bool {
				return len(k) == len(prefix)+12 && bytes.HasPrefix(k, prefix)
			}

			k, _ := cur.Seek(prefix)
			if !ok(k) {
				break
			}

			if orderASC && afterPkey != nil {
				for {
					if bytes.Compare(k[len(k)-12:], afterPkey) <= 0 {
						k, _ = cur.Next()
						if !ok(k) {
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
					} else if len(k) != len(prefix)+12 || !bytes.HasPrefix(k, prefix) || (afterPkey != nil && bytes.Compare(k[len(k)-12:], afterPkey) >= 0) {
						k, _ = cur.Prev()
						break
					}
				}
			}

			for {
				if !ok(k) {
					break
				}
				pkey := k[len(k)-12:]
				if pkey != nil {
					if i == 0 {
						data := metaBucket.Get(pkey)
						if data != nil {
							post, err := q.PostFromBytes(data)
							if err == nil && bytes.Equal(post.PKey.Bytes(), pkey) &&
								(filter == nil || filter(post) == true) &&
								(!res.HasStatus || post.Status == res.Status) {
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
							if bytes.Equal(p.PKey.Bytes(), pkey) {
								if i == pl-1 {
									posts = append(posts, *p)
								} else {
									a[i] = append(a[i], p)
								}
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

	if len(res.Keys) > 0 {
		for _, post := range posts {
			postkvBucket := kvBucket.Bucket([]byte(post.ID))
			if postkvBucket != nil {
				if res.KeysHasWildcard {
					c := postkvBucket.Cursor()
					for k, v := c.First(); k != nil; k, v = c.Next() {
						post.KV[string(k)] = v
					}
				} else {
					for key := range res.Keys {
						if kl := len(key); kl > 1 && strings.HasSuffix(key, "*") {
							c := postkvBucket.Cursor()
							prefix := []byte(key[:kl-1])
							for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
								post.KV[string(k)] = v
							}
						} else {
							v := postkvBucket.Get([]byte(key))
							if v != nil {
								post.KV[key] = v
							}
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
		q.Resolve(&res)
	}

	var metaData []byte
	if len(res.ID) > 0 {
		uidIndexBucket := tx.t.Bucket(postindexKey).Bucket(postuidKey)
		pkey := uidIndexBucket.Get([]byte(res.ID))
		if pkey != nil {
			metaData = tx.t.Bucket(postmetaKey).Get(pkey)
		}
	}
	if metaData == nil {
		return nil, ErrNotFound
	}

	post, err := q.PostFromBytes(metaData)
	if err != nil {
		return nil, err
	}

	if len(res.Keys) > 0 {
		postkvBucket := tx.t.Bucket(postkvKey).Bucket([]byte(post.ID))
		if postkvBucket != nil {
			if res.KeysHasWildcard {
				c := postkvBucket.Cursor()
				for k, v := c.First(); k != nil; k, v = c.Next() {
					post.KV[string(k)] = v
				}
			} else {
				for key := range res.Keys {
					if kl := len(key); kl > 1 && strings.HasSuffix(key, "*") {
						c := postkvBucket.Cursor()
						prefix := []byte(key[:kl-1])
						for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
							post.KV[string(k)] = v
						}
					} else {
						v := postkvBucket.Get([]byte(key))
						if v != nil {
							post.KV[key] = v
						}
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
	uidIndexBucket := indexBucket.Bucket(postuidKey)

newpost:
	post := q.NewPost()
	for _, q := range qs {
		q.Apply(post)
	}

	// ensure the post.ID is unique
	if uidIndexBucket.Get([]byte(post.ID)) != nil {
		goto newpost
	}

	err := metaBucket.Put(post.PKey.Bytes(), post.MetaData())
	if err != nil {
		return nil, err
	}

	err = uidIndexBucket.Put([]byte(post.ID), post.PKey.Bytes())
	if err != nil {
		return nil, err
	}

	if len(post.Alias) > 0 {
		if uidIndexBucket.Get([]byte(post.Alias)) != nil {
			return nil, ErrDuplicateAlias
		}
		err = uidIndexBucket.Put([]byte(post.Alias), post.PKey.Bytes())
		if err != nil {
			return nil, err
		}
	}

	if len(post.Owner) > 0 {
		ownerIndexBucket := indexBucket.Bucket(postownerKey)
		keypath := [][]byte{[]byte(post.Owner), post.PKey.Bytes()}
		err = ownerIndexBucket.Put(bytes.Join(keypath, []byte{0}), []byte{1})
		if err != nil {
			return nil, err
		}
	}

	if len(post.Tags) > 0 {
		tagIndexBucket := indexBucket.Bucket(posttagKey)
		for _, tag := range post.Tags {
			keypath := [][]byte{[]byte(tag), post.PKey.Bytes()}
			err = tagIndexBucket.Put(bytes.Join(keypath, []byte{0}), []byte{1})
			if err != nil {
				return nil, err
			}
		}
	}

	postkvBucket, err := kvBucket.CreateBucketIfNotExists([]byte(post.ID))
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
	uidIndexBucket := indexBucket.Bucket(postuidKey)

	post, err := tx.Get(qs...)
	if err != nil {
		return nil, err
	}

	copy := post.Clone()
	for _, q := range qs {
		q.Apply(copy)
	}

	shouldUpdateMeta := copy.Status != post.Status

	// update alias index
	if copy.Alias != post.Alias {
		if uidIndexBucket.Get([]byte(copy.Alias)) != nil {
			return nil, ErrDuplicateAlias
		}
		if len(post.Alias) > 0 {
			err = uidIndexBucket.Delete([]byte(post.Alias))
			if err != nil {
				return nil, err
			}
		}
		if len(copy.Alias) > 0 {
			err = uidIndexBucket.Put([]byte(copy.Alias), copy.PKey.Bytes())
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
			keypath := [][]byte{[]byte(post.Owner), post.PKey.Bytes()}
			err = ownerIndexBucket.Delete(bytes.Join(keypath, []byte{0}))
			if err != nil {
				return nil, err
			}
		}
		if len(copy.Owner) > 0 {
			keypath := [][]byte{[]byte(copy.Owner), copy.PKey.Bytes()}
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
				keypath := [][]byte{[]byte(tag), post.PKey.Bytes()}
				err = tagIndexBucket.Delete(bytes.Join(keypath, []byte{0}))
				if err != nil {
					return nil, err
				}
			}
		}
		if len(copy.Tags) > 0 {
			for _, tag := range copy.Tags {
				keypath := [][]byte{[]byte(tag), copy.PKey.Bytes()}
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
		err = metaBucket.Put(copy.PKey.Bytes(), copy.MetaData())
		if err != nil {
			return nil, err
		}
	}

	if len(copy.KV) > 0 {
		postkvBucket := kvBucket.Bucket([]byte(copy.ID))
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
func (tx *Tx) DeleteKV(qs ...q.Query) (n int, err error) {
	post, err := tx.Get(qs...)
	if err != nil {
		return
	}

	var res q.Resolver
	for _, q := range qs {
		q.Resolve(&res)
	}

	if len(res.Keys) > 0 {
		kvBucket := tx.t.Bucket(postkvKey)
		postkvBucket := kvBucket.Bucket([]byte(post.ID))
		if postkvBucket != nil {
			for key := range res.Keys {
				err = postkvBucket.Delete([]byte(key))
				if err != nil {
					return
				}
				n++
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
	uidIndexBucket := indexBucket.Bucket(postuidKey)
	posts := tx.List(qs...)

	for _, post := range posts {
		err = metaBucket.Delete(post.PKey.Bytes())
		if err != nil {
			return
		}

		err = uidIndexBucket.Delete([]byte(post.ID))
		if err != nil {
			return
		}

		if len(post.Alias) > 0 {
			err = uidIndexBucket.Delete([]byte(post.Alias))
			if err != nil {
				return
			}
		}

		if len(post.Owner) > 0 {
			ownerIndexBucket := indexBucket.Bucket(postownerKey)
			keypath := [][]byte{[]byte(post.Owner), []byte(post.ID)}
			err = ownerIndexBucket.Delete(bytes.Join(keypath, []byte{0}))
			if err != nil {
				return
			}
		}

		if len(post.Tags) > 0 {
			tagIndexBucket := indexBucket.Bucket(posttagKey)
			for _, tag := range post.Tags {
				keypath := [][]byte{[]byte(tag), []byte(post.ID)}
				err = tagIndexBucket.Delete(bytes.Join(keypath, []byte{0}))
				if err != nil {
					return
				}
			}
		}

		err = kvBucket.DeleteBucket([]byte(post.ID))
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
