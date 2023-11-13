package post

import (
	"encoding/binary"
	"errors"
	"time"

	"github.com/rs/xid"
)

const (
	postMetaVersion = 1
	postMetaPrefix  = "POST"
)

var (
	// ErrBadPostMeta specifies the error of bad post meta
	ErrBadPostMeta = errors.New("bad post meta")
)

// A Post specifies a post of postdb.
type Post struct {
	PKey    [12]byte
	Owner   uint32
	Crtime  uint32
	Modtime uint32
	Status  uint8
	Alias   string
	Tags    []string
	KV      map[string][]byte
	// todo: Rank uint32
	// todo: Revision uint32
}

// New returns a new post.
func New() *Post {
	now := uint32(time.Now().Unix())
	post := &Post{
		PKey:    xid.New(),
		Status:  1,
		Crtime:  now,
		Modtime: now,
		Tags:    []string{},
		KV:      map[string][]byte{},
	}
	return post
}

// FromBytes parses a post from bytes.
func FromBytes(data []byte) (post *Post, err error) {
	dl := len(data)
	if dl < 6 {
		return nil, ErrBadPostMeta
	}

	for i, c := range postMetaPrefix {
		if data[i] != byte(c) {
			return nil, ErrBadPostMeta
		}
	}

	var checksum byte
	for i := 0; i < dl-1; i++ {
		checksum += data[i]
	}
	if data[dl-1] != checksum {
		return nil, ErrBadPostMeta
	}

	version := data[4]
	switch version {
	case 1:
		return decodeV1(data[5 : dl-1])
	default:
		return nil, ErrBadPostMeta
	}
}

func decodeV1(data []byte) (*Post, error) {
	dl := len(data)
	if dl < 33-6 {
		return nil, ErrBadPostMeta
	}

	var i int
	var pkey xid.ID
	copy(pkey[:], data[i:i+12])
	i += 12
	owner := binary.BigEndian.Uint32(data[i : i+8])
	i += 4
	crtime := binary.BigEndian.Uint32(data[i : i+8])
	i += 4
	modtime := binary.BigEndian.Uint32(data[i : i+8])
	i += 4
	status := data[i]
	i++
	aliasLen := int(data[i])
	i++
	tagN := int(data[i])
	i++
	if i+aliasLen > dl {
		return nil, ErrBadPostMeta
	}
	alias := data[i : i+aliasLen]
	i += aliasLen
	tags := make([]string, tagN)
	for t := 0; t < tagN; t++ {
		tl := int(data[i])
		tEnd := i + 1 + tl
		if tEnd > dl {
			return nil, ErrBadPostMeta
		}
		tags[t] = string(data[i+1 : tEnd])
		i += 1 + tl
	}

	return &Post{
		PKey:    pkey,
		Owner:   owner,
		Crtime:  crtime,
		Modtime: modtime,
		Status:  uint8(status),
		Alias:   string(alias),
		Tags:    tags,
		KV:      map[string][]byte{},
	}, nil
}

// Bytes returns the bytes of the post structure.
// data structure:
// "POST"(4) | version(1) | pkey(12) | owner(4) | crtime(4) | modtime(4)| status(1) | aliasLen(1) | tagsLen(1) | alias(aliasLen) | tags([1+tagLen]*N) | checksum(1)
func (p *Post) Bytes() []byte {
	aliasLen := len(p.Alias)
	dataLen := 33 + aliasLen
	for _, tag := range p.Tags {
		dataLen += 1 + len(tag)
	}
	buf := make([]byte, dataLen)
	i := 0
	copy(buf, postMetaPrefix)
	i += 4
	buf[i] = postMetaVersion
	i++
	copy(buf[i:], p.PKey[:])
	i += 12
	binary.BigEndian.PutUint32(buf[i:], p.Owner)
	i += 4
	binary.BigEndian.PutUint32(buf[i:], p.Crtime)
	i += 4
	binary.BigEndian.PutUint32(buf[i:], p.Modtime)
	i += 4
	buf[i] = byte(p.Status)
	i++
	buf[i] = byte(aliasLen)
	i++
	buf[i] = byte(len(p.Tags))
	i++
	if aliasLen > 0 {
		copy(buf[i:], []byte(p.Alias))
		i += aliasLen
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
	for j := 0; j < dataLen-1; j++ {
		checksum += buf[j]
	}
	buf[i] = checksum
	return buf
}

func (p *Post) ID() string {
	return xid.ID(p.PKey).String()
}

// Clone clones the post
func (p *Post) Clone() *Post {
	clone := &Post{
		PKey:    p.PKey,
		Alias:   p.Alias,
		Owner:   p.Owner,
		Status:  p.Status,
		Crtime:  p.Crtime,
		Modtime: p.Modtime,
		Tags:    make([]string, len(p.Tags)),
		KV:      map[string][]byte{},
	}
	copy(clone.Tags, p.Tags)
	for k, v := range p.KV {
		b := make([]byte, len(v))
		copy(b, v)
		clone.KV[k] = b
	}
	return clone
}
