package q

import (
	"encoding/binary"
	"time"

	"github.com/rs/xid"
)

// A Post specifies a post of postdb.
type Post struct {
	PKey    xid.ID
	ID      string
	Alias   string
	Status  uint8
	Owner   string
	Crtime  uint32
	Modtime uint32
	Tags    []string
	KV      KV
}

// NewPost returns a new post.
func NewPost() *Post {
	now := uint32(time.Now().Unix())
	post := &Post{
		PKey:    xid.New(),
		ID:      NewID(),
		Status:  1,
		Crtime:  now,
		Modtime: now,
		Tags:    []string{},
		KV:      KV{},
	}
	return post
}

// PostFromBytes parses a post from bytes.
func PostFromBytes(data []byte) (post *Post, err error) {
	dl := len(data)
	// v1 min 50
	if dl < 50 {
		return nil, errPostMeta
	}

	for i, c := range postPrefix {
		if data[i] != c {
			return nil, errPostMeta
		}
	}

	var checksum byte
	for i := 0; i < dl-1; i++ {
		checksum += data[i]
	}
	if data[dl-1] != checksum {
		return nil, errPostMeta
	}

	version := data[4]
	switch version {
	case 1:
		return decodeV1(data[5 : dl-1])
	default:
		return nil, errPostMeta
	}
}

func decodeV1(data []byte) (*Post, error) {
	dl := len(data)
	var i int
	var pkey xid.ID
	copy(pkey[:], data[i:i+12])
	i += 12
	id := string(data[i : i+20])
	i += 20
	crtime := binary.BigEndian.Uint32(data[i : i+8])
	i += 4
	modtime := binary.BigEndian.Uint32(data[i : i+8])
	i += 4
	status := data[i]
	i++
	aliasLen := int(data[i])
	i++
	ownerLen := int(data[i])
	i++
	tagN := int(data[i])
	i++
	if i+aliasLen > dl {
		return nil, errPostMeta
	}
	alias := data[i : i+aliasLen]
	i += aliasLen
	if i+ownerLen > dl {
		return nil, errPostMeta
	}
	owner := data[i : i+ownerLen]
	i += ownerLen
	tags := make([]string, tagN)
	for t := 0; t < tagN; t++ {
		tl := int(data[i])
		tEnd := i + 1 + tl
		if tEnd > dl {
			return nil, errPostMeta
		}
		tags[t] = string(data[i+1 : tEnd])
		i += 1 + tl
	}

	return &Post{
		PKey:    pkey,
		ID:      id,
		Alias:   string(alias),
		Owner:   string(owner),
		Status:  uint8(status),
		Crtime:  crtime,
		Modtime: modtime,
		Tags:    tags,
		KV:      KV{},
	}, nil
}

// Clone clones the post
func (p *Post) Clone(qs ...Query) *Post {
	clone := &Post{
		PKey:    p.PKey,
		ID:      p.ID,
		Alias:   p.Alias,
		Owner:   p.Owner,
		Status:  p.Status,
		Crtime:  p.Crtime,
		Modtime: p.Modtime,
		Tags:    make([]string, len(p.Tags)),
		KV:      KV{},
	}
	for i, t := range p.Tags {
		clone.Tags[i] = t
	}
	for k, v := range p.KV {
		b := make([]byte, len(v))
		copy(b, v)
		clone.KV[k] = b
	}
	return clone
}

// MetaData returns the meta data of post.
// data structure:
// "POST"(4) | version(1) | pkey(12) | id(20) | crtime(4) | modtime(4)| status(1) | aliasLen(1) | ownerLen(1) | tagsN(1) | alias(aliasLen) | owner(ownerLen) | tags([1+tagLen]*tagN) | checksum(1)
func (p *Post) MetaData() []byte {
	aliasLen := len(p.Alias)
	ownerLen := len(p.Owner)
	metaLen := 50 + aliasLen + ownerLen
	for _, tag := range p.Tags {
		metaLen += 1 + len(tag)
	}
	buf := make([]byte, metaLen)
	i := 0
	copy(buf, postPrefix)
	i += 4
	buf[i] = PostMetaDataVersion
	i++
	copy(buf[i:], p.PKey.Bytes())
	i += 12
	copy(buf[i:], []byte(p.ID))
	i += 20
	binary.BigEndian.PutUint32(buf[i:], p.Crtime)
	i += 4
	binary.BigEndian.PutUint32(buf[i:], p.Modtime)
	i += 4
	buf[i] = byte(p.Status)
	i++
	buf[i] = byte(aliasLen)
	i++
	buf[i] = byte(ownerLen)
	i++
	buf[i] = byte(len(p.Tags))
	i++
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
		if tl > 0 {
			copy(buf[i+1:], []byte(tag))
		}
		i += 1 + tl
	}
	var checksum byte
	for j := 0; j < metaLen-1; j++ {
		checksum += buf[j]
	}
	buf[i] = checksum
	return buf
}
