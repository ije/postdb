package q

// A KV map
type KV map[string][]byte

// Has reports whether h has the provided key defined
func (kv KV) Has(key string) bool {
	_, ok := kv[key]
	return ok
}

// Get gets the first value associated with the given key.
func (kv KV) Get(key string) []byte {
	return kv[key]
}

// Set sets the header entries associated with key to the
// single element value.
func (kv KV) Set(key string, value []byte) {
	kv[key] = value
}

// Del deletes the values associated with key.
func (kv KV) Del(key string) {
	delete(kv, key)
}

// Error implements the Query interface
func (kv KV) Error() error {
	return nil
}
