package postdb

import (
	"errors"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrDuplicateSlug = errors.New("duplicate slug")
)
