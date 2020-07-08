package postdb

import (
	"crypto/tls"

	"github.com/ije/puddle"
)

type Client struct {
	pool *puddle.Pool
}

type ConnConfig struct {
	Host      string
	Port      uint16
	Secret    string
	TLSConfig *tls.Config
}

func Connect(conn *ConnConfig) (*Client, error) {
	return nil, nil
}
