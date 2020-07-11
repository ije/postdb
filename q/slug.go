package q

type slugQuery string

// Slug returns a slug Query
func Slug(slug string) Query {
	return slugQuery(slug)
}

// QueryType implements the Query interface
func (q slugQuery) QueryType() string {
	return "slug"
}
