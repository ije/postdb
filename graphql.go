package postdb

import (
	"net/http"

	"github.com/graphql-go/graphql"
)

type GraphQLMux struct {
	db     Database
	schema *graphql.Schema
}

func Graphql(db Database) *GraphQLMux {
	return nil
}

func (mux *GraphQLMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
