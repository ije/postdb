package q

import (
	"encoding/binary"
	"errors"
	"time"
)

var (
	postPrefix  = []byte{'P', 'O', 'S', 'T'}
	errPostMeta = errors.New("bad post meta data")
)

// A Post specifies a post of postdb.
type Post struct {
	ID     ObjectID
	Alias  string
	Status uint8
	Owner  string
	Crtime uint64
	Tags   []string
	KV     KV
}

// NewPost returns a new post.
func NewPost() *Post {
	return &Post{
		ID:     NewID(),
		Status: 1,
		Crtime: uint64(time.Now().UnixNano() / 1e6),
		Tags:   []string{},
		KV:     KV{},
	}
}

// PostFromBytes parses a post from bytes.
func PostFromBytes(data []byte) (*Post, error) {
	dl := len(data)
	if dl < 29 {
		return nil, errPostMeta
	}

	for i, l := 0, len(postPrefix); i < l; i++ {
		if data[i] != postPrefix[i] {
			return nil, errPostMeta
		}
	}

	var v byte
	for j := 0; j < dl-1; j++ {
		v += data[j]
	}
	if data[dl-1] != v {
		return nil, errPostMeta
	}

	var id ObjectID
	copy(id[:], data[4:16])

	status := data[16]
	crtime := binary.BigEndian.Uint64(data[17:25])
	aliasLen := int(data[25])
	ownerLen := int(data[26])
	tagN := int(data[27])
	i := 28
	alias := data[i : i+aliasLen]
	i += aliasLen
	owner := data[i : i+ownerLen]
	i += ownerLen
	tags := make([]string, tagN)
	for t := 0; t < tagN; t++ {
		tl := int(data[i])
		tEnd := i + 1 + tl
		if tEnd >= dl {
			return nil, errPostMeta
		}
		tags[t] = string(data[i+1 : tEnd])
		i += 1 + tl
	}

	return &Post{
		ID:     id,
		Alias:  string(alias),
		Owner:  string(owner),
		Status: uint8(status),
		Crtime: crtime,
		Tags:   tags,
		KV:     KV{},
	}, nil
}

// Clone clones the post.
func (p *Post) Clone(qs ...Query) *Post {
	copy := &Post{
		ID:     p.ID,
		Alias:  p.Alias,
		Owner:  p.Owner,
		Status: p.Status,
		Crtime: p.Crtime,
		Tags:   make([]string, len(p.Tags)),
		KV:     KV{},
	}
	for i, t := range p.Tags {
		copy.Tags[i] = t
	}
	for k, v := range p.KV {
		copy.KV[k] = v
	}
	return copy
}

// MetaData returns the meta data of post.
// data structure:
// "POST"(4) | id(12) | status(1) | crtime(8) | aliasLen(1) | ownerLen(1) | tagsN(1) | alias(aliasLen) | owner(ownerLen) | tags([1+tagLen]*tagN) | checksum(1)
func (p *Post) MetaData() []byte {
	aliasLen := len(p.Alias)
	ownerLen := len(p.Owner)
	metaLen := 29 + aliasLen + ownerLen
	for _, tag := range p.Tags {
		metaLen += 1 + len(tag)
	}
	buf := make([]byte, metaLen)
	copy(buf, postPrefix)
	copy(buf[4:], p.ID.Bytes())
	buf[16] = byte(p.Status)
	binary.BigEndian.PutUint64(buf[17:], p.Crtime)
	buf[25] = byte(aliasLen)
	buf[26] = byte(ownerLen)
	buf[27] = byte(len(p.Tags))
	i := 28
	if aliasLen > 0 {
		copy(buf[i:], []byte(p.Alias))
		i += aliasLen
	}
	if ownerLen > 0 {
		copy(buf[i:], []byte(p.Owner))
		i += ownerLen
	}
	for _, tag := range p.Tags {
		tl := len(tag)
		buf[i] = byte(tl)
		copy(buf[i+1:], []byte(tag))
		i += 1 + tl
	}
	var v byte
	for j := 0; j < metaLen-1; j++ {
		v += buf[j]
	}
	buf[i] = v
	return buf
}

// ApplyQuery applies a query.
func (p *Post) ApplyQuery(query Query) {
	switch q := query.(type) {
	case aliasQuery:
		p.Alias = string(q)

	case ownerQuery:
		p.Owner = string(q)

	case statusQuery:
		p.Status = uint8(q)

	case tagsQuery:
		p.Tags = q

	case KV:
		if p.KV == nil {
			p.KV = KV{}
		}
		for k, v := range q {
			if len(k) > 0 && v != nil {
				p.KV[k] = v
			}
		}
	}
}
