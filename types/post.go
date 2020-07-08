package types

// A Post specifies a post of postdb
type Post struct {
	ID             string            `json:"id"`
	Parent         string            `json:"parent,omitempty"`
	Link           string            `json:"link,omitempty"`
	Type           string            `json:"type"`
	Tags           []string          `json:"tags"`
	CurrentVersion uint32            `json:"currentVersion"`
	Status         PostStatus        `json:"status"`
	Crtime         int64             `json:"crtime"`
	Content        map[string][]byte `json:"-"`
}
