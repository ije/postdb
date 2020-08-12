package postdb

import (
	"crypto/tls"
	"strings"

	"github.com/ije/puddle"
)

type ConnConfig struct {
	Host      string
	Port      uint16
	User      string
	Password  string
	TLSConfig *tls.Config
}

type Client struct {
	pool *puddle.Pool
}

func Connect(config ConnConfig) (*Client, error) {
	return nil, nil
}

// DB opens a client database
func (c *Client) DB(name string) (*ClientDB, error) {
	return &ClientDB{
		db:     strings.TrimSpace(strings.ToLower(name)),
		client: c,
	}, nil
}
