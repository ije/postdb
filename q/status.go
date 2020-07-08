package q

const (
	DRAFT QueryStatus = iota
	NORMAL
	DELETED
)

// A QueryStatus specifies the status of Query
type QueryStatus uint8

// Status returns a status Query
func Status(status QueryStatus, extra ...[]QueryStatus) Query {
	return nil
}
