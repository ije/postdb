package q

import "fmt"

const (
	DRAFT PostStatus = iota
	NORMAL
	DELETED
)

// A PostStatus specifies the status of Query
type PostStatus uint8

func (status PostStatus) String() string {
	switch status {
	case DRAFT:
		return "DRAFT"
	case NORMAL:
		return "NORMAL"
	case DELETED:
		return "DELETED"
	default:
		return fmt.Sprintf("STATUS_%d", status)
	}
}

// Status returns a status Query
func Status(status PostStatus, extra ...[]PostStatus) Query {
	return nil
}
