package postdb

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/ije/postdb/internal/post"
	"github.com/ije/postdb/internal/util"
	"github.com/ije/postdb/q"
	bolt "go.etcd.io/bbolt"
)

// Tx represents a transaction on the database.
type Tx struct {
	ns []byte
	tx *bolt.Tx
}

// List returns some posts
func (tx *Tx) List(qs ...q.Query) (posts []post.Post) {
	metaBucket := tx.bucket(keyPostMeta)
	indexBucket := tx.bucket(keyPostIndex)
	aliasIndexBucket := indexBucket.Bucket(keyPostAlias)

	var indexCur *bolt.Cursor
	var prefixs [][]byte
	var filter func(*post.Post) bool
	var n uint32
	var res q.Resolver

	for _, q := range qs {
		q.Resolve(&res)
	}
	queryOwner := res.Owner > 0
	orderASC := res.Order >= q.ASC

	/* query post list by the IDs */
	if len(res.IDs) > 0 {
		posts = make([]post.Post, len(res.IDs))
		for i, id := range res.IDs {
			v := metaBucket.Get(id[:])
			if v != nil {
				post, err := post.FromBytes(v)
				if err == nil && bytes.Equal(post.PKey[:], id[:]) && (len(res.Tags) == 0 || util.Contains(post.Tags, res.Tags)) && (!queryOwner || post.Owner == res.Owner) {
					if len(res.Keys) > 0 {
						tx.readKV(post, res.Keys)
					}
					if res.Filter(*post.Clone()) {
						posts[i] = *post
						n++
					}
				}
			}
		}
		posts = posts[:n]
		return
	}

	/* query post list by the Alias */
	if len(res.Alias) > 0 {
		posts = make([]post.Post, len(res.Alias))
		for i, alias := range res.Alias {
			pkey := aliasIndexBucket.Get([]byte(alias))
			if pkey != nil {
				v := metaBucket.Get(pkey)
				if v != nil {
					post, err := post.FromBytes(v)
					if err == nil && bytes.Equal(post.PKey[:], pkey) && (len(res.Tags) == 0 || util.Contains(post.Tags, res.Tags)) && (!queryOwner || post.Owner == res.Owner) {
						if len(res.Keys) > 0 {
							tx.readKV(post, res.Keys)
						}
						if res.Filter(*post.Clone()) {
							posts[i] = *post
							n++
						}
					}
				}
			}
		}
		posts = posts[:n]
		return
	}

	if len(res.Tags) > 0 {
		indexCur = indexBucket.Bucket(keyPostTag).Cursor()
		prefixs = make([][]byte, len(res.Tags))
		for i, tag := range res.Tags {
			prefixs[i] = util.ToPrefix(tag)
		}
		if queryOwner {
			filter = func(post *post.Post) bool {
				return post.Owner == res.Owner
			}
		}
	} else if queryOwner {
		indexCur = indexBucket.Bucket(keyPostOwner).Cursor()
		prefixs = [][]byte{util.Join(util.Uint32ToBytes(res.Owner), nil, 0)}
	}

	if indexCur != nil {
		pl := len(prefixs)
		a := make([][]*post.Post, pl)
		for i, prefix := range prefixs {
			k, _ := indexCur.Seek(prefix)
			ok := func(k []byte) bool {
				return len(k) == len(prefix)+12 && bytes.HasPrefix(k, prefix)
			}
			if !ok(k) {
				break
			}

			// move cursor to the last key firstly when order by DESC
			if !orderASC {
				for {
					k, _ = indexCur.Next()
					if k == nil {
						k, _ = indexCur.Last()
						break
					} else if !ok(k) {
						k, _ = indexCur.Prev()
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
							post, err := post.FromBytes(data)
							if err == nil && bytes.Equal(post.PKey[:], pkey) && (filter == nil || filter(post) == true) {
								if len(res.Keys) > 0 {
									tx.readKV(post, res.Keys)
								}
								if res.Filter(*post.Clone()) {
									if pl == 1 {
										if res.Offset == 0 || n >= res.Offset {
											posts = append(posts, *post)
										}
										n++
										if res.Limit > 0 && uint32(len(posts)) >= res.Limit {
											return
										}
									} else {
										a[0] = append(a[0], post)
									}
								}
							}
						}
					} else if past := a[i-1]; len(past) > 0 {
						for _, p := range past {
							if bytes.Equal(p.PKey[:], pkey) {
								if i == pl-1 {
									if res.Offset == 0 || n >= res.Offset {
										posts = append(posts, *p)
									}
									n++
									if res.Limit > 0 && uint32(len(posts)) >= res.Limit {
										return
									}
								} else {
									a[i] = append(a[i], p)
								}
							}
						}
					}
				}
				if orderASC {
					k, _ = indexCur.Next()
				} else {
					k, _ = indexCur.Prev()
				}
			}
		}
	} else {
		c := metaBucket.Cursor()
		var k, v []byte
		if orderASC {
			k, v = c.First()
		} else {
			k, v = c.Last()
		}
		for {
			if k == nil {
				break
			}
			post, err := post.FromBytes(v)
			if err == nil && bytes.Equal(post.PKey[:], k) {
				if len(res.Keys) > 0 {
					tx.readKV(post, res.Keys)
				}
				if res.Filter(*post.Clone()) {
					if res.Offset == 0 || n >= res.Offset {
						posts = append(posts, *post)
					}
					n++
					if res.Limit > 0 && uint32(len(posts)) >= res.Limit {
						return
					}
				}
			}
			if orderASC {
				k, v = c.Next()
			} else {
				k, v = c.Prev()
			}
		}
	}
	return
}

func (tx *Tx) readKV(post *post.Post, keys []string) {
	postkvBucket := tx.bucket(keyPostKV).Bucket(post.PKey[:])
	if postkvBucket == nil {
		return
	}

	var rest []string
	for _, key := range keys {
		kl := len(key)
		if kl > 0 {
			if strings.HasSuffix(key, "*") {
				c := postkvBucket.Cursor()
				// key equals "*"
				if kl == 1 {
					for k, v := c.First(); k != nil; k, v = c.Next() {
						post.KV[string(k)] = v
					}
					return
				}
				prefix := []byte(key[:kl-1])
				for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
					post.KV[string(k)] = v
				}
			} else if _, ok := post.KV[key]; !ok {
				rest = append(rest, key)
			}
		}
	}
	for _, key := range rest {
		v := postkvBucket.Get([]byte(key))
		if v != nil {
			post.KV[key] = v
		}
	}
}

// Get returns the post
func (tx *Tx) Get(qs ...q.Query) (*post.Post, error) {
	var res q.Resolver
	for _, q := range qs {
		q.Resolve(&res)
	}

	var metadata []byte
	if len(res.IDs) > 0 {
		pkey := res.IDs[0][:]
		metadata = tx.bucket(keyPostMeta).Get(pkey)
	} else if len(res.Alias) > 0 {
		pkey := tx.bucket(keyPostIndex).Bucket(keyPostAlias).Get([]byte(res.Alias[0]))
		if pkey != nil {
			metadata = tx.bucket(keyPostMeta).Get(pkey)
		}
	}
	if metadata == nil {
		return nil, ErrNotFound
	}

	post, err := post.FromBytes(metadata)
	if err != nil {
		return nil, err
	}

	if len(res.Keys) > 0 {
		tx.readKV(post, res.Keys)
	}

	return post, nil
}

// Put puts a new post
func (tx *Tx) Put(qs ...q.Query) (*post.Post, error) {
	post := post.New()
	for _, q := range qs {
		q.Apply(post)
	}
	err := tx.PutPost(post)
	if err != nil {
		return nil, err
	}
	return post, nil
}

// PutPost puts a new post
func (tx *Tx) PutPost(post *post.Post) (err error) {
	metaBucket := tx.bucket(keyPostMeta)
	indexBucket := tx.bucket(keyPostIndex)
	kvBucket := tx.bucket(keyPostKV)
	aliasIndexBucket := indexBucket.Bucket(keyPostAlias)

	if metaBucket.Get(post.PKey[:]) != nil {
		return fmt.Errorf("duplicate pkey %v", post.PKey)
	}
	err = metaBucket.Put(post.PKey[:], post.Bytes())
	if err != nil {
		return
	}

	if len(post.Alias) > 0 {
		if aliasIndexBucket.Get([]byte(post.Alias)) != nil {
			return ErrDuplicateAlias
		}
		err = aliasIndexBucket.Put([]byte(post.Alias), post.PKey[:])
		if err != nil {
			return
		}
	}

	if post.Owner > 0 {
		ownerIndexBucket := indexBucket.Bucket(keyPostOwner)
		keypath := util.Join(util.Uint32ToBytes(post.Owner), post.PKey[:], 0)
		err = ownerIndexBucket.Put(keypath, []byte{1})
		if err != nil {
			return
		}
	}

	if len(post.Tags) > 0 {
		tagIndexBucket := indexBucket.Bucket(keyPostTag)
		for _, tag := range post.Tags {
			keypath := util.Join([]byte(tag), post.PKey[:], 0)
			err = tagIndexBucket.Put(keypath, []byte{1})
			if err != nil {
				return
			}
		}
	}

	postkvBucket, err := kvBucket.CreateBucketIfNotExists(post.PKey[:])
	if err != nil {
		return
	}

	if len(post.KV) > 0 {
		for k, v := range post.KV {
			err = postkvBucket.Put([]byte(k), v)
			if err != nil {
				return
			}
		}
	}

	return
}

// Update updates the post
func (tx *Tx) Update(qs ...q.Query) (ok bool, err error) {
	metaBucket := tx.bucket(keyPostMeta)
	indexBucket := tx.bucket(keyPostIndex)
	kvBucket := tx.bucket(keyPostKV)
	idIndexBucket := indexBucket.Bucket(keyPostAlias)

	post, err := tx.Get(qs...)
	if err != nil {
		if err == ErrNotFound {
			err = nil
		}
		return
	}

	copy := post.Clone()
	for _, q := range qs {
		q.Apply(copy)
	}

	shouldUpdateMeta := copy.Status != post.Status

	// update alias index
	if copy.Alias != post.Alias {
		if idIndexBucket.Get([]byte(copy.Alias)) != nil {
			err = ErrDuplicateAlias
			return
		}
		if len(post.Alias) > 0 {
			err = idIndexBucket.Delete([]byte(post.Alias))
			if err != nil {
				return
			}
		}
		if len(copy.Alias) > 0 {
			err = idIndexBucket.Put([]byte(copy.Alias), copy.PKey[:])
			if err != nil {
				return
			}
		}
		shouldUpdateMeta = true
	}

	// update owner index
	if copy.Owner != post.Owner {
		ownerIndexBucket := indexBucket.Bucket(keyPostOwner)
		if post.Owner > 0 {
			keypath := util.Join(util.Uint32ToBytes(post.Owner), post.PKey[:], 0)
			err = ownerIndexBucket.Delete(keypath)
			if err != nil {
				return
			}
		}
		if copy.Owner > 0 {
			keypath := util.Join(util.Uint32ToBytes(copy.Owner), copy.PKey[:], 0)
			err = ownerIndexBucket.Put(keypath, []byte{1})
			if err != nil {
				return
			}
		}
		shouldUpdateMeta = true
	}

	// update tags index
	if strings.Join(copy.Tags, "") != strings.Join(post.Tags, "") {
		tagIndexBucket := indexBucket.Bucket(keyPostTag)
		if len(post.Tags) > 0 {
			for _, tag := range post.Tags {
				keypath := util.Join([]byte(tag), post.PKey[:], 0)
				err = tagIndexBucket.Delete(keypath)
				if err != nil {
					return
				}
			}
		}
		if len(copy.Tags) > 0 {
			for _, tag := range copy.Tags {
				keypath := util.Join([]byte(tag), copy.PKey[:], 0)
				err = tagIndexBucket.Put(keypath, []byte{1})
				if err != nil {
					return
				}
			}
		}
		shouldUpdateMeta = true
	}

	if len(copy.KV) > 0 {
		postkvBucket := kvBucket.Bucket(copy.PKey[:])
		for k, v := range copy.KV {
			err = postkvBucket.Put([]byte(k), v)
			if err != nil {
				return
			}
		}
		shouldUpdateMeta = true
	}

	if shouldUpdateMeta {
		copy.Modtime = uint32(time.Now().Unix())
		err = metaBucket.Put(copy.PKey[:], copy.Bytes())
		if err != nil {
			return
		}
	}

	ok = true
	return
}

// DeleteKV deletes the post kv
func (tx *Tx) DeleteKV(qs ...q.Query) error {
	post, err := tx.Get(qs...)
	if err != nil {
		return err
	}

	var res q.Resolver
	for _, q := range qs {
		q.Resolve(&res)
	}

	if len(res.Keys) > 0 {
		kvBucket := tx.bucket(keyPostKV)
		postkvBucket := kvBucket.Bucket(post.PKey[:])
		if postkvBucket != nil {
			var rest []string
			for _, key := range res.Keys {
				kl := len(key)
				if kl > 0 {
					if strings.HasSuffix(key, "*") {
						c := postkvBucket.Cursor()
						// key equals "*"
						if kl == 1 {
							for k, _ := c.First(); k != nil; k, _ = c.Next() {
								err := postkvBucket.Delete(k)
								if err != nil {
									return err
								}
							}
							return nil
						}
						prefix := []byte(key[:kl-1])
						for k, _ := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, _ = c.Next() {
							err := postkvBucket.Delete(k)
							if err != nil {
								return err
							}
						}
					} else if _, ok := post.KV[key]; !ok {
						rest = append(rest, key)
					}
				}
			}
			for _, key := range rest {
				err := postkvBucket.Delete([]byte(key))
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// Delete deletes the post
func (tx *Tx) Delete(qs ...q.Query) (n int, err error) {
	if len(qs) == 0 {
		return 0, nil
	}

	metaBucket := tx.bucket(keyPostMeta)
	indexBucket := tx.bucket(keyPostIndex)
	kvBucket := tx.bucket(keyPostKV)
	aliasIndexBucket := indexBucket.Bucket(keyPostAlias)
	posts := tx.List(qs...)

	for _, post := range posts {
		err = metaBucket.Delete(post.PKey[:])
		if err != nil {
			return
		}

		if len(post.Alias) > 0 {
			err = aliasIndexBucket.Delete([]byte(post.Alias))
			if err != nil {
				return
			}
		}

		if post.Owner > 0 {
			ownerIndexBucket := indexBucket.Bucket(keyPostOwner)
			keypath := util.Join(util.Uint32ToBytes(post.Owner), post.PKey[:], 0)
			err = ownerIndexBucket.Delete(keypath)
			if err != nil {
				return
			}
		}

		if len(post.Tags) > 0 {
			tagIndexBucket := indexBucket.Bucket(keyPostTag)
			for _, tag := range post.Tags {
				keypath := util.Join([]byte(tag), post.PKey[:], 0)
				err = tagIndexBucket.Delete(keypath)
				if err != nil {
					return
				}
			}
		}

		err = kvBucket.DeleteBucket(post.PKey[:])
		if err == bolt.ErrBucketNotFound {
			err = nil
		}
		if err != nil {
			return
		}
	}

	n = len(posts)
	return
}

func (tx *Tx) bucket(name []byte) *bolt.Bucket {
	if len(tx.ns) > 0 {
		return tx.tx.Bucket(util.Join(tx.ns, name, 0))
	}
	return tx.tx.Bucket(name)
}

// Rollback closes the transaction and ignores all previous updates. Read-only
// transactions must be rolled back and not committed.
func (tx *Tx) Rollback() error {
	return tx.tx.Rollback()
}

// Commit writes all changes to disk and updates the meta page.
// Returns an error if a disk write error occurs, or if Commit is
// called on a read-only transaction.
func (tx *Tx) Commit() error {
	return tx.tx.Commit()
}

// WriteTo writes the entire database to a writer.
// If err == nil then exactly tx.Size() bytes will be written into the writer.
func (tx *Tx) WriteTo(w io.Writer) (int64, error) {
	return tx.tx.WriteTo(w)
}
