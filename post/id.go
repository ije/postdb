package post

import (
	"crypto/rand"
)

const (
	idLen   = 20
	idChars = "abcdefghijklmnopqrstuv0123456789"
)

// NewID returns a new ID string in base32
func NewID() string {
	r := make([]byte, idLen)
	buf := make([]byte, idLen)
	rand.Read(r)
	buf[0] = idChars[r[0]%22] // always start with char [a-v]
	for i := 1; i < idLen; i++ {
		buf[i] = idChars[r[i]%32]
	}
	return string(buf)
}
