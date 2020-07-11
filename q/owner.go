package q

type ownerQuery string

// Owner returns a owner Query
func Owner(name string) Query {
	return ownerQuery(name)
}

// QueryType implements the Query interface
func (q ownerQuery) QueryType() string {
	return "owner"
}
