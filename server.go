package postdb

import (
	"crypto/tls"
)

type Server struct {
	DB        *NSDB
	Port      uint16
	Secret    string
	TLSConfig *tls.Config
}

func (s *Server) Serve() error {
	return nil
}
