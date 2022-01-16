package q

import "github.com/ije/postdb/post"

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

// Delete deletes the values associated with key.
func (kv KV) Delete(key string) {
	delete(kv, key)
}

// Apply implements the Query interface
func (kv KV) Apply(p *post.Post) {
	if p.KV == nil {
		p.KV = KV{}
	}
	for k, v := range kv {
		if len(k) > 0 && v != nil {
			p.KV[k] = v
		}
	}
}

// Resolve implements the Query interface
func (kv KV) Resolve(r *Resolver) {}
