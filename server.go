package postdb

import (
	"crypto/tls"
)

type ServerConfig struct {
	Port      uint16
	Secret    string
	TLSConfig *tls.Config
}

func ListenAndServe(db NSDB, config *ServerConfig) error {
	return nil
}
