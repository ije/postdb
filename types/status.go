package types

import (
	"fmt"
)

const (
	DRAFT PostStatus = iota
	NORMAL
	DELETED
)

// A PostStatus specifies the status of Post
type PostStatus uint8

func (status PostStatus) String() string {
	switch status {
	case DRAFT:
		return "draft"
	case NORMAL:
		return "normal"
	case DELETED:
		return "deleted"
	default:
		return fmt.Sprintf("status_%d", status)
	}
}
