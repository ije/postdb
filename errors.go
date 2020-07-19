package postdb

import (
	"errors"
)

var (
	ErrNotFound       = errors.New("not found")
	ErrDuplicateAlias = errors.New("duplicate alias")
)
