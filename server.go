package postdb

import (
	"crypto/tls"
)

type ServerConfig struct {
	Port      uint16
	Secret    string
	TLSConfig *tls.Config
}

type Server struct {
	db Database
}

func ListenAndServe(db Database, config *ServerConfig) error {
	return nil
}
