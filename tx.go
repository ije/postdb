package postdb

import (
	"bytes"
	"fmt"
	"io"

	"github.com/postui/postdb/q"
	"github.com/rs/xid"
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

func (tx *Tx) GetPosts(qs ...q.Query) (posts []q.Post, err error) {
	metaBucket := tx.t.Bucket(postmetaKey)
	indexBucket := tx.t.Bucket(postindexKey)
	kvBucket := tx.t.Bucket(postkvKey)

	var ret q.QueryResult
	var cur *bolt.Cursor
	var prefixs [][]byte
	var n int

	for _, q := range qs {
		ret.ApplyQuery(q)
	}

	queryTags := len(ret.Tags) > 0
	queryType := len(ret.Type) > 0
	queryOwner := len(ret.Owner) > 0

	if queryTags {
		cur = indexBucket.Bucket(tagsKey).Cursor()
		prefixs = make([][]byte, len(ret.Tags))
		for i, tag := range ret.Tags {
			prefix := make([]byte, len(tag)+1)
			copy(prefix, []byte(tag))
			if queryType {
				keypath := bytes.Join([][]byte{[]byte(tag), []byte(ret.Type)}, []byte{0})
				prefix = make([]byte, len(keypath)+1)
				copy(prefix, keypath)
			}
			prefixs[i] = prefix
		}
	} else if queryOwner {
		prefix := make([]byte, len(ret.Owner)+1)
		copy(prefix, []byte(ret.Owner))
		if queryType {
			keypath := bytes.Join([][]byte{[]byte(ret.Owner), []byte(ret.Type)}, []byte{0})
			prefix = make([]byte, len(keypath)+1)
			copy(prefix, keypath)
		}
		cur = indexBucket.Bucket(ownersKey).Cursor()
		prefixs = [][]byte{prefix}
	} else if queryType {
		prefix := make([]byte, len(ret.Type)+1)
		copy(prefix, []byte(ret.Type))
		cur = indexBucket.Bucket(typesKey).Cursor()
		prefixs = [][]byte{prefix}
	} else {
		c := metaBucket.Cursor()
		var k, v []byte
		var post *q.Post
		if len(ret.Aftar) == 12 {
			k, v = c.Seek(ret.Aftar)
		} else {
			k, v = c.First()
		}
		for ; len(k) == 12; k, v = c.Next() {
			post, err = q.ParsePostMeta(v)
			if err != nil {
				posts = nil
				return
			}
			posts = append(posts, *post)
			n++
			if ret.Limit > 0 && n >= ret.Limit {
				break
			}
		}
	}

	if cur != nil && len(prefixs) > 0 {
		var k []byte
		var p *q.Post
		for _, prefix := range prefixs {
			if len(ret.Aftar) == 12 {
				k, _ = cur.Seek(bytes.Join([][]byte{prefix, ret.Aftar}, []byte{}))
			} else {
				k, _ = cur.Seek(prefix)
			}
			for ; len(k) > 12+len(prefix) && bytes.HasPrefix(k, prefix); k, _ = cur.Next() {
				id := k[len(k)-12:]
				data := metaBucket.Get(id)
				if data != nil {
					p, err = q.ParsePostMeta(data)
					if err != nil {
						posts = nil
						return
					}
					if queryTags && queryOwner && p.Owner != ret.Owner {
						continue
					}
					posts = append(posts, *p)
					n++
					if ret.Limit > 0 && n >= ret.Limit {
						break
					}
				}
			}
		}
	}

	if len(ret.Keys) > 0 {
		for _, post := range posts {
			postkvBucket := kvBucket.Bucket(post.ID)
			for _, key := range ret.Keys {
				v := postkvBucket.Get([]byte(key))
				if v != nil {
					post.KV[key] = v
				}
			}
		}
	}
	return
}

func (tx *Tx) GetPost(idOrSlug string, keys q.Keys) (*q.Post, error) {
	metaBucket := tx.t.Bucket(postmetaKey)
	var postMetaData []byte
	if len(idOrSlug) == 20 {
		id, err := xid.FromString(idOrSlug)
		if err == nil {
			postMetaData = metaBucket.Get(id.Bytes())
		}
	}
	if postMetaData == nil {
		slugsBucket := tx.t.Bucket(postindexKey).Bucket(slugsKey)
		id := slugsBucket.Get([]byte(idOrSlug))
		if id != nil {
			postMetaData = metaBucket.Get(id)
		}
	}
	if postMetaData == nil {
		return nil, ErrNotFound
	}

	post, err := q.ParsePostMeta(postMetaData)
	if err != nil {
		return nil, err
	}

	if len(keys) > 0 {
		kvBucket := tx.t.Bucket(postkvKey).Bucket(post.ID)
		for _, key := range keys {
			v := kvBucket.Get([]byte(key))
			if v != nil {
				post.KV[key] = v
			}
		}
	}

	return post, nil
}

func (tx *Tx) AddPost(postType string, qs ...q.Query) (*q.Post, error) {
	post := q.NewPost(postType)
	for _, q := range qs {
		post.ApplyQuery(q)
	}

	metaBucket := tx.t.Bucket(postmetaKey)
	indexBucket := tx.t.Bucket(postindexKey)
	kvBucket := tx.t.Bucket(postkvKey)

	err := metaBucket.Put(post.ID, post.MetaData())
	if err != nil {
		return nil, err
	}

	if len(post.KV) > 0 {
		postkvBucket, err := kvBucket.CreateBucketIfNotExists(post.ID)
		if err != nil {
			return nil, err
		}
		for k, v := range post.KV {
			err = postkvBucket.Put([]byte(k), v)
			if err != nil {
				return nil, err
			}
		}
	}

	if len(post.Slug) > 0 {
		slugsBucket := indexBucket.Bucket(slugsKey)
		err = slugsBucket.Put([]byte(post.Slug), post.ID)
		if err != nil {
			return nil, err
		}
	}

	hasType := len(post.Type) > 0
	if hasType {
		typesBucket := indexBucket.Bucket(typesKey)
		keypath := [][]byte{[]byte(post.Type), post.ID}
		err = typesBucket.Put(bytes.Join(keypath, []byte{0}), []byte{1})
		if err != nil {
			return nil, err
		}
	}

	if len(post.Owner) > 0 {
		ownersBucket := indexBucket.Bucket(ownersKey)
		var keypath [][]byte
		if hasType {
			keypath = [][]byte{[]byte(post.Owner), []byte(post.Type), post.ID}
		} else {
			keypath = [][]byte{[]byte(post.Owner), post.ID}
		}
		err = ownersBucket.Put(bytes.Join(keypath, []byte{0}), []byte{1})
		if err != nil {
			return nil, err
		}
	}

	if len(post.Tags) > 0 {
		tagsBucket := indexBucket.Bucket(tagsKey)
		for _, tag := range post.Tags {
			var keypath [][]byte
			if hasType {
				keypath = [][]byte{[]byte(tag), []byte(post.Type), post.ID}
			} else {
				keypath = [][]byte{[]byte(tag), post.ID}
			}
			err = tagsBucket.Put(bytes.Join(keypath, []byte{0}), []byte{1})
			if err != nil {
				return nil, err
			}
		}
	}

	return post, nil
}

func (tx *Tx) UpdatePost(idOrSlug string, qs ...q.Query) error {
	post, err := tx.GetPost(idOrSlug, nil)
	if err != nil {
		return err
	}

	copy := post.Clone()
	for _, q := range qs {
		copy.ApplyQuery(q)
	}

	// todo: update db
	return nil
}

func (tx *Tx) RemovePost(idOrSlug string) error {
	post, err := tx.GetPost(idOrSlug, nil)
	if err != nil {
		return err
	}

	metaBucket := tx.t.Bucket(postmetaKey)
	indexBucket := tx.t.Bucket(postindexKey)
	kvBucket := tx.t.Bucket(postkvKey)

	err = metaBucket.Delete(post.ID)
	if err != nil {
		return err
	}

	if len(post.Slug) > 0 {
		slugsBucket := indexBucket.Bucket(slugsKey)
		err = slugsBucket.Delete([]byte(post.Slug))
		if err != nil {
			return err
		}
		fmt.Printf("rmslug: slug=%s", post.Slug)
	}

	if len(post.Type) > 0 {
		typesBucket := indexBucket.Bucket(typesKey)
		keypath := [][]byte{[]byte(post.Type), post.ID}
		err = typesBucket.Delete(bytes.Join(keypath, []byte{0}))
		if err != nil {
			return err
		}
		fmt.Printf("rmtype: keypath=%s", keypath)
	}

	if len(post.Owner) > 0 {
		ownersBucket := indexBucket.Bucket(ownersKey)
		var keypath [][]byte
		if len(post.Type) > 0 {
			keypath = [][]byte{[]byte(post.Owner), []byte(post.Type), post.ID}
		} else {
			keypath = [][]byte{[]byte(post.Owner), post.ID}
		}
		err = ownersBucket.Delete(bytes.Join(keypath, []byte{0}))
		if err != nil {
			return err
		}
		fmt.Printf("rmowner: keypath=%s", keypath)
	}

	if len(post.Tags) > 0 {
		tagsBucket := indexBucket.Bucket(tagsKey)
		for _, tag := range post.Tags {
			var keypath [][]byte
			if len(post.Type) > 0 {
				keypath = [][]byte{[]byte(tag), []byte(post.Type), post.ID}
			} else {
				keypath = [][]byte{[]byte(tag), post.ID}
			}
			err = tagsBucket.Delete(bytes.Join(keypath, []byte{0}))
			if err != nil {
				return err
			}
			fmt.Printf("rmtag: keypath=%s", bytes.Join(keypath, []byte{' '}))
		}
	}

	err = kvBucket.DeleteBucket(post.ID)
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
