package postdb

import (
	"net/http"

	"github.com/graphql-go/graphql"
)

type GraphQL struct {
	Schema *graphql.Schema
}

func NewGraphql(db NSDB) *GraphQL {
	return nil
}

func (gql *GraphQL) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
