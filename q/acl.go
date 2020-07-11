package q

const (
	PRIVATE ACL = iota
	PUBLIC_READ
	PUBLIC_READ_WRITE
)

// A ACL specifies the acl of Post
type ACL uint8

func (status ACL) String() string {
	switch status {
	case PRIVATE:
		return "private"
	case PUBLIC_READ:
		return "public-read"
	case PUBLIC_READ_WRITE:
		return "public-read-write"
	default:
		return "unkonwn"
	}
}

// QueryType implements the Query interface
func (status ACL) QueryType() string {
	return "acl"
}
