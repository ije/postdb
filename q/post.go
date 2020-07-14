package q

import (
	"encoding/binary"
	"errors"
	"time"

	"github.com/rs/xid"
)

var (
	postPrefix  = []byte{'P', 'O', 'S', 'T'}
	errPostMeta = errors.New("bad post meta data")
)

// A Post specifies a post of postdb.
type Post struct {
	ID     xid.ID
	Slug   string
	Type   string
	Status uint8
	Owner  string
	Crtime uint64
	Tags   []string
	KV     KV
}

// NewPost returns a new post.
func NewPost() *Post {
	return &Post{
		ID:     xid.New(),
		Status: 1,
		Crtime: uint64(time.Now().UnixNano() / 1e6),
		Tags:   []string{},
		KV:     KV{},
	}
}

// PostFromBytes parses a post from bytes.
func PostFromBytes(data []byte) (*Post, error) {
	dl := len(data)
	if dl < 30 {
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

	id, err := xid.FromBytes(data[4:16])
	if err != nil {
		return nil, errPostMeta
	}

	status := data[16]
	crtime := binary.BigEndian.Uint64(data[17:25])
	slugLen := int(data[25])
	typeLen := int(data[26])
	ownerLen := int(data[27])
	tagN := int(data[28])
	i := 29
	slug := data[i : i+slugLen]
	i += slugLen
	postType := data[i : i+typeLen]
	i += typeLen
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
		Slug:   string(slug),
		Type:   string(postType),
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
		Type:   p.Type,
		Slug:   p.Slug,
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
// "POST"(4) | id(12) | status(1) | crtime(8) | slugLen(1) | typeLen(1) | ownerLen(1) | tagsN(1) | slug(slugLen) | type(typeLen) | owner(ownerLen) | tags([1+tagLen]*tagN) | checksum(1)
func (p *Post) MetaData() []byte {
	slugLen := len(p.Slug)
	typeLen := len(p.Type)
	ownerLen := len(p.Owner)
	metaLen := 30 + slugLen + typeLen + ownerLen
	for _, tag := range p.Tags {
		metaLen += 1 + len(tag)
	}
	buf := make([]byte, metaLen)
	copy(buf, postPrefix)
	copy(buf[4:], p.ID.Bytes())
	buf[16] = byte(p.Status)
	binary.BigEndian.PutUint64(buf[17:], p.Crtime)
	buf[25] = byte(slugLen)
	buf[26] = byte(typeLen)
	buf[27] = byte(ownerLen)
	buf[28] = byte(len(p.Tags))
	i := 29
	if slugLen > 0 {
		copy(buf[i:], []byte(p.Slug))
		i += slugLen
	}
	if typeLen > 0 {
		copy(buf[i:], []byte(p.Type))
		i += typeLen
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
	case slugQuery:
		p.Slug = toLowerTrim(string(q))

	case typeQuery:
		p.Type = toLowerTrim(string(q))

	case statusQuery:
		p.Status = uint8(q)

	case ownerQuery:
		p.Owner = toLowerTrim(string(q))

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
