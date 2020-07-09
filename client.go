package postdb

import (
	"crypto/tls"

	"github.com/ije/puddle"
)

type ConnConfig struct {
	Host      string
	Port      uint16
	Secret    string
	TLSConfig *tls.Config
}

type Client struct {
	pool *puddle.Pool
}

func Connect(conn *ConnConfig) (*Client, error) {
	return nil, nil
}
