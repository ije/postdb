package q

type typeQuery string

// Type returns a type Query
func Type(t string) Query {
	return typeQuery(t)
}

// QueryType implements the Query interface
func (q typeQuery) QueryType() string {
	return "type"
}
