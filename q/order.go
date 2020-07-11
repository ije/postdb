package q

const (
	ASC Order = iota
	DESC
)

type Order uint8

// QueryType implements the Query interface
func (o Order) QueryType() string {
	return "order"
}
