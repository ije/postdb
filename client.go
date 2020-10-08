package postdb

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"strings"

	"github.com/ije/puddle"
)

type ConnConfig struct {
	Host        string
	Port        uint16
	User        string
	Password    string
	MaxPoolSize int32
	TLSConfig   *tls.Config
}

type Client struct {
	config ConnConfig
	pool   *puddle.Pool
}

func Connect(config ConnConfig) (*Client, error) {
	constructor := func(context.Context) (interface{}, error) {
		return net.Dial("tcp", fmt.Sprintf("%s:%d", config.Host, config.Port))
	}
	destructor := func(value interface{}) {
		value.(net.Conn).Close()
	}
	pool := puddle.NewPool(constructor, destructor, config.MaxPoolSize)
	return &Client{config, pool}, nil
}

// DB opens a client database
func (c *Client) DB(name string) (*ClientDB, error) {
	return &ClientDB{
		db:     strings.TrimSpace(strings.ToLower(name)),
		client: c,
	}, nil
}
