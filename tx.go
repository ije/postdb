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

func (tx *Tx) GetPosts(qs ...q.Query) (posts []q.Post, err error) {
	metaBucket := tx.t.Bucket(postmetaKey)
	indexBucket := tx.t.Bucket(postindexKey)
	kvBucket := tx.t.Bucket(postkvKey)

	var cur *bolt.Cursor
	var prefixs [][]byte
	var n int

	var res q.Resolver
	for _, q := range qs {
		res.Apply(q)
	}

	queryTags := len(res.Tags) > 0
	queryType := len(res.Type) > 0
	queryOwner := len(res.Owner) > 0

	if queryTags {
		cur = indexBucket.Bucket(tagsKey).Cursor()
		prefixs = make([][]byte, len(res.Tags))
		for i, tag := range res.Tags {
			prefix := make([]byte, len(tag)+1)
			copy(prefix, []byte(tag))
			if queryType {
				keypath := bytes.Join([][]byte{[]byte(tag), []byte(res.Type)}, []byte{0})
				prefix = make([]byte, len(keypath)+1)
				copy(prefix, keypath)
			}
			prefixs[i] = prefix
		}
	} else if queryOwner {
		prefix := make([]byte, len(res.Owner)+1)
		copy(prefix, []byte(res.Owner))
		if queryType {
			keypath := bytes.Join([][]byte{[]byte(res.Owner), []byte(res.Type)}, []byte{0})
			prefix = make([]byte, len(keypath)+1)
			copy(prefix, keypath)
		}
		cur = indexBucket.Bucket(ownersKey).Cursor()
		prefixs = [][]byte{prefix}
	} else if queryType {
		prefix := make([]byte, len(res.Type)+1)
		copy(prefix, []byte(res.Type))
		cur = indexBucket.Bucket(typesKey).Cursor()
		prefixs = [][]byte{prefix}
	} else {
		c := metaBucket.Cursor()
		var k, v []byte
		var post *q.Post
		if len(res.Aftar) == 12 {
			k, v = c.Seek(res.Aftar)
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
			if res.Limit > 0 && n >= res.Limit {
				break
			}
		}
	}

	if cur != nil && len(prefixs) > 0 {
		var k []byte
		var p *q.Post
		for _, prefix := range prefixs {
			if len(res.Aftar) == 12 {
				k, _ = cur.Seek(bytes.Join([][]byte{prefix, res.Aftar}, []byte{}))
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
					if queryTags && queryOwner && p.Owner != res.Owner {
						continue
					}
					posts = append(posts, *p)
					n++
					if res.Limit > 0 && n >= res.Limit {
						break
					}
				}
			}
		}
	}

	if len(res.Keys) > 0 {
		for _, post := range posts {
			postkvBucket := kvBucket.Bucket(post.ID)
			for _, key := range res.Keys {
				v := postkvBucket.Get([]byte(key))
				if v != nil {
					post.KV[key] = v
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

	metaBucket := tx.t.Bucket(postmetaKey)
	var postMetaData []byte
	if len(res.ID) == 12 {
		postMetaData = metaBucket.Get(res.ID)
	} else if len(res.Slug) > 0 {
		slugsBucket := tx.t.Bucket(postindexKey).Bucket(slugsKey)
		id := slugsBucket.Get([]byte(res.Slug))
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

	if len(res.Keys) > 0 {
		kvBucket := tx.t.Bucket(postkvKey).Bucket(post.ID)
		for _, key := range res.Keys {
			v := kvBucket.Get([]byte(key))
			if v != nil {
				post.KV[key] = v
			}
		}
	}

	return post, nil
}

func (tx *Tx) AddPost(qs ...q.Query) (*q.Post, error) {
	post := q.NewPost("")
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

func (tx *Tx) UpdatePost(qs ...q.Query) error {
	post, err := tx.GetPost(qs...)
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

func (tx *Tx) RemovePost(qs ...q.Query) error {
	post, err := tx.GetPost(qs...)
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
