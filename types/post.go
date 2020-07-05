package types

const (
	POST_STAT_Draft PostStatus = iota
	POST_STAT_NORMAL
	POST_STAT_DELETED
)

// A PostStatus specifies the status of User
type PostStatus uint8

func (status PostStatus) String() string {
	switch status {
	case POST_STAT_Draft:
		return "draft"
	case POST_STAT_NORMAL:
		return "normal"
	case POST_STAT_DELETED:
		return "deleted"
	default:
		return ""
	}
}

// A Post specifies the post of postdb
type Post struct {
	ID             string            `json:"id"`
	Parent         string            `json:"parent,omitempty"`
	Link           string            `json:"link,omitempty"`
	Type           string            `json:"type"`
	Category       string            `json:"category"`
	Tags           []string          `json:"tags"`
	CurrentVersion uint32            `json:"currentVersion"`
	Status         PostStatus        `json:"status"`
	Crtime         int64             `json:"crtime"`
	Content        map[string]string `json:"content"`
}
