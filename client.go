package postdb

import (
	"crypto/tls"
	"fmt"
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
	name = strings.TrimSpace(strings.ToLower(name))
	if !nameReg.MatchString(name) {
		return nil, fmt.Errorf("invalid name '%s'", name)
	}
	return &ClientDB{db: name, client: c}, nil
}
