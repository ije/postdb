package postdb

import (
	"net/http"

	"github.com/graphql-go/graphql"
)

type GraphQL struct {
	db     Database
	schema *graphql.Schema
}

func NewGraphql(db Database) *GraphQL {
	return nil
}

func (gql *GraphQL) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
