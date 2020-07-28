package q

import (
	"errors"
)

var (
	postPrefix  = []byte{'P', 'O', 'S', 'T'}
	errPostMeta = errors.New("bad post meta data")
)

const (
	// DESC specifies the order of DESC
	DESC uint8 = iota
	// ASC specifies the order of ASC
	ASC
)

const (
	// PostMetaDataVersion specifies the current version of a post meta data
	PostMetaDataVersion = 1
)
