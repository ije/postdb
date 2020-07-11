package q

type rangeQuery struct {
	after string
	limit int
}

// Range returns a range Query
func Range(after string, limit int) Query {
	return rangeQuery{
		after: after,
		limit: limit,
	}
}

// QueryType implements the Query interface
func (q rangeQuery) QueryType() string {
	return "range"
}
