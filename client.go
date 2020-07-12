package postdb

import (
	"crypto/tls"
	"io"

	"github.com/ije/puddle"
	"github.com/postui/postdb/q"
)

type ConnConfig struct {
	Host      string
	Port      uint16
	Secret    string
	Namespace string
	TLSConfig *tls.Config
}

type Client struct {
	pool *puddle.Pool
}

func Connect(config ConnConfig) (*Client, error) {
	return nil, nil
}

// Begin starts a new transaction.
func (c *Client) Begin(writable bool) (*Tx, error) {
	return nil, nil
}

// GetValue returns the value for a key in the database.
func (c *Client) GetValue(key string) ([]byte, error) {
	return nil, nil
}

// PutValue sets the value for a key in the database.
func (c *Client) PutValue(key string, value []byte) error {
	return nil
}

func (c *Client) GetPost(qs ...q.Query) (*q.Post, error) {
	return nil, nil
}

func (c *Client) GetPosts(qs ...q.Query) ([]q.Post, error) {
	return nil, nil
}

func (c *Client) AddPost(qs ...q.Query) (*q.Post, error) {
	return nil, nil
}

func (c *Client) UpdatePost(qs ...q.Query) error {
	return nil
}

func (c *Client) RemovePost(qs ...q.Query) error {
	return nil
}

// WriteTo writes the entire database to a writer.
func (c *Client) WriteTo(w io.Writer) (int64, error) {
	return 0, nil
}
