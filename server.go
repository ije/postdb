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

type Server struct {
	db *NSDB
}

func ListenAndServe(config ServerConfig) error {
	return nil
}
