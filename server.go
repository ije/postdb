package postdb

import (
	"crypto/tls"
)

type ServerConfig struct {
	DB        *NSDB
	Port      uint16
	Secret    string
	TLSConfig *tls.Config
}

type server struct {
	db *NSDB
}

func Serve(config ServerConfig) error {
	return nil
}
