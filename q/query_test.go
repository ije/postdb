package q

import (
	"strings"
	"testing"

	"github.com/rs/xid"
)

func TestQuery(t *testing.T) {
	var ret QueryResult
	for _, q := range []Query{
		Owner("admin"),
		Type("news"),
		Keys{"title", "content"},
		Tags("hello", "world"),
		Range(xid.New().String(), 100),
		DESC,
	} {
		ret.ApplyQuery(q)
	}

	toBe(t, "Owner", ret.Owner, "admin")
	toBe(t, "Type", ret.Type, "news")
	toBe(t, "Tags", strings.Join(ret.Tags, " "), "hello world")
	toBe(t, "Keys", strings.Join(ret.Keys, " "), "title content")
	toBe(t, "Aftar", len(ret.Aftar), 12)
	toBe(t, "Limit", ret.Limit, 100)
	toBe(t, "Order", ret.Order, DESC)

	t.Log(ret)
}
