package postdb

import (
	"bytes"
	"fmt"
	"io"

	"github.com/postui/postdb/q"
	bolt "go.etcd.io/bbolt"
)

// Tx represents a transaction on the database.
type Tx struct {
	t *bolt.Tx
}

func (tx *Tx) GetValue(key string) []byte {
	return tx.t.Bucket(valuesKey).Get([]byte(key))
}

func (tx *Tx) PutValue(key string, value []byte) error {
	return tx.t.Bucket(valuesKey).Put([]byte(key), value)
}

func (tx *Tx) GetPosts(qs ...q.Query) (posts []q.Post) {
	var cur *bolt.Cursor
	var prefixs [][]byte
	var filter func(*q.Post) bool
	var n uint32

	var res q.Resolver
	for _, q := range qs {
		res.Apply(q)
	}

	queryType := len(res.Type) > 0
	queryTags := len(res.Tags) > 0
	queryOwner := len(res.Owner) > 0
	queryAfter := len(res.After) == 12
	orderASC := res.Order >= q.ASC

	metaBucket := tx.t.Bucket(postmetaKey)
	indexBucket := tx.t.Bucket(postindexKey)
	kvBucket := tx.t.Bucket(postkvKey)

	if queryTags {
		cur = indexBucket.Bucket(tagKey).Cursor()
		prefixs = make([][]byte, len(res.Tags))
		for i, tag := range res.Tags {
			p := make([]byte, len(tag)+1)
			copy(p, []byte(tag))
			if queryType {
				keypath := bytes.Join([][]byte{[]byte(tag), []byte(res.Type)}, []byte{0})
				p = make([]byte, len(keypath)+1)
				copy(p, keypath)
			}
			prefixs[i] = p
		}
		if queryOwner {
			filter = func(post *q.Post) bool {
				return post.Owner == res.Owner
			}
		}
	} else if queryOwner {
		prefix := make([]byte, len(res.Owner)+1)
		copy(prefix, []byte(res.Owner))
		if queryType {
			keypath := bytes.Join([][]byte{[]byte(res.Owner), []byte(res.Type)}, []byte{0})
			prefix = make([]byte, len(keypath)+1)
			copy(prefix, keypath)
		}
		cur = indexBucket.Bucket(ownerKey).Cursor()
		prefixs = [][]byte{prefix}
	} else if queryType {
		prefix := make([]byte, len(res.Type)+1)
		copy(prefix, []byte(res.Type))
		cur = indexBucket.Bucket(typeKey).Cursor()
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
						if err == nil && bytes.Compare(post.ID.Bytes(), id) == 0 && (filter == nil || filter(post) == true) {
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

	if len(res.Keys) > 0 {
		for _, post := range posts {
			postkvBucket := kvBucket.Bucket(post.ID.Bytes())
			if postkvBucket != nil {
				if res.WildcardKey {
					c := postkvBucket.Cursor()
					for k, v := c.First(); k != nil; k, v = c.Next() {
						post.KV[string(k)] = v
					}
				} else {
					for _, key := range res.Keys {
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

func (tx *Tx) GetPost(qs ...q.Query) (*q.Post, error) {
	var res q.Resolver
	for _, q := range qs {
		res.Apply(q)
	}

	var metaData []byte
	if len(res.ID) == 12 {
		metaData = tx.t.Bucket(postmetaKey).Get(res.ID)
	} else if len(res.Slug) > 0 {
		slugIndexBucket := tx.t.Bucket(postindexKey).Bucket(slugKey)
		id := slugIndexBucket.Get([]byte(res.Slug))
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

	if len(res.Keys) > 0 {
		postkvBucket := tx.t.Bucket(postkvKey).Bucket(post.ID.Bytes())
		if postkvBucket != nil {
			if res.WildcardKey {
				c := postkvBucket.Cursor()
				for k, v := c.First(); k != nil; k, v = c.Next() {
					post.KV[string(k)] = v
				}
			} else {
				for _, key := range res.Keys {
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

func (tx *Tx) AddPost(qs ...q.Query) (*q.Post, error) {
	post := q.NewPost()
	for _, q := range qs {
		post.ApplyQuery(q)
	}

	err := tx.t.Bucket(postmetaKey).Put(post.ID.Bytes(), post.MetaData())
	if err != nil {
		return nil, err
	}

	postkvBucket, err := tx.t.Bucket(postkvKey).CreateBucketIfNotExists(post.ID.Bytes())
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

	hasType := len(post.Type) > 0
	indexBucket := tx.t.Bucket(postindexKey)

	if len(post.Slug) > 0 {
		slugIndexBucket := indexBucket.Bucket(slugKey)
		if slugIndexBucket.Get([]byte(post.Slug)) != nil {
			return nil, ErrDuplicateSlug
		}
		err = slugIndexBucket.Put([]byte(post.Slug), post.ID.Bytes())
		if err != nil {
			return nil, err
		}
	}

	if hasType {
		typeIndexBucket := indexBucket.Bucket(typeKey)
		keypath := [][]byte{[]byte(post.Type), post.ID.Bytes()}
		err = typeIndexBucket.Put(bytes.Join(keypath, []byte{0}), []byte{1})
		if err != nil {
			return nil, err
		}
	}

	if len(post.Owner) > 0 {
		ownerIndexBucket := indexBucket.Bucket(ownerKey)
		var keypath [][]byte
		if hasType {
			keypath = [][]byte{[]byte(post.Owner), []byte(post.Type), post.ID.Bytes()}
		} else {
			keypath = [][]byte{[]byte(post.Owner), post.ID.Bytes()}
		}
		err = ownerIndexBucket.Put(bytes.Join(keypath, []byte{0}), []byte{1})
		if err != nil {
			return nil, err
		}
	}

	if len(post.Tags) > 0 {
		tagIndexBucket := indexBucket.Bucket(tagKey)
		for _, tag := range post.Tags {
			var keypath [][]byte
			if hasType {
				keypath = [][]byte{[]byte(tag), []byte(post.Type), post.ID.Bytes()}
			} else {
				keypath = [][]byte{[]byte(tag), post.ID.Bytes()}
			}
			err = tagIndexBucket.Put(bytes.Join(keypath, []byte{0}), []byte{1})
			if err != nil {
				return nil, err
			}
		}
	}

	return post, nil
}

func (tx *Tx) UpdatePost(qs ...q.Query) error {
	post, err := tx.GetPost(qs...)
	if err != nil {
		return err
	}

	copy := post.Clone()
	for _, q := range qs {
		copy.ApplyQuery(q)
	}

	metaData := post.MetaData()
	copyMetaData := copy.MetaData()
	if !bytes.Equal(metaData, copyMetaData) {
		metaBucket := tx.t.Bucket(postmetaKey)
		err := metaBucket.Put(copy.ID.Bytes(), copyMetaData)
		if err != nil {
			return err
		}

		// todo: update indexs
	}

	if len(post.KV) > 0 {
		postkvBucket := tx.t.Bucket(postkvKey).Bucket(post.ID.Bytes())
		for k, v := range post.KV {
			err = postkvBucket.Put([]byte(k), v)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (tx *Tx) RemovePost(qs ...q.Query) error {
	post, err := tx.GetPost(qs...)
	if err != nil {
		return err
	}

	metaBucket := tx.t.Bucket(postmetaKey)
	indexBucket := tx.t.Bucket(postindexKey)
	kvBucket := tx.t.Bucket(postkvKey)

	err = metaBucket.Delete(post.ID.Bytes())
	if err != nil {
		return err
	}

	if len(post.Slug) > 0 {
		slugIndexBucket := indexBucket.Bucket(slugKey)
		err = slugIndexBucket.Delete([]byte(post.Slug))
		if err != nil {
			return err
		}
		fmt.Printf("rmslug: slug=%s", post.Slug)
	}

	if len(post.Type) > 0 {
		typeIndexBucket := indexBucket.Bucket(typeKey)
		keypath := [][]byte{[]byte(post.Type), post.ID.Bytes()}
		err = typeIndexBucket.Delete(bytes.Join(keypath, []byte{0}))
		if err != nil {
			return err
		}
		fmt.Printf("rmtype: keypath=%s", keypath)
	}

	if len(post.Owner) > 0 {
		ownerIndexBucket := indexBucket.Bucket(ownerKey)
		var keypath [][]byte
		if len(post.Type) > 0 {
			keypath = [][]byte{[]byte(post.Owner), []byte(post.Type), post.ID.Bytes()}
		} else {
			keypath = [][]byte{[]byte(post.Owner), post.ID.Bytes()}
		}
		err = ownerIndexBucket.Delete(bytes.Join(keypath, []byte{0}))
		if err != nil {
			return err
		}
		fmt.Printf("rmowner: keypath=%s", keypath)
	}

	if len(post.Tags) > 0 {
		tagIndexBucket := indexBucket.Bucket(tagKey)
		for _, tag := range post.Tags {
			var keypath [][]byte
			if len(post.Type) > 0 {
				keypath = [][]byte{[]byte(tag), []byte(post.Type), post.ID.Bytes()}
			} else {
				keypath = [][]byte{[]byte(tag), post.ID.Bytes()}
			}
			err = tagIndexBucket.Delete(bytes.Join(keypath, []byte{0}))
			if err != nil {
				return err
			}
			fmt.Printf("rmtag: keypath=%s", bytes.Join(keypath, []byte{' '}))
		}
	}

	err = kvBucket.DeleteBucket(post.ID.Bytes())
	return err
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
