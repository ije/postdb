package q

type KV map[string][]byte
type Keys []string

// QueryType implements the Query interface
func (kv KV) QueryType() string {
	return "kv"
}

// QueryType implements the Query interface
func (keys Keys) QueryType() string {
	return "keys"
}
