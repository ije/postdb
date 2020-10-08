package postdb

import (
	"crypto/tls"
	"fmt"
	"net"
	"sync"
	"time"
)

type ServerConfig struct {
	DBPath string
	Port   uint16
}

type Server struct {
	dbPoolLock sync.RWMutex
	dbPool     map[string]*DB
	coreDB     *DB
	ServerConfig
}

func (s *Server) CreateUser(username string, password string) error {
	return nil
}

func (s *Server) CreateDB(name string, owner string) error {
	return nil
}

func (s *Server) Serve() error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", s.Port))
	if err != nil {
		return err
	}

	return s.serve(l)
}

func (s *Server) ServeTLS(certFile string, keyFile string) error {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return err
	}

	config := &tls.Config{Certificates: []tls.Certificate{cert}}
	l, err := tls.Listen("tcp", fmt.Sprintf(":%d", s.Port), config)
	if err != nil {
		return err
	}

	return s.serve(l)
}

func (s *Server) serve(l net.Listener) error {
	var tempDelay time.Duration // how long to sleep on accept failure

	for {
		conn, err := l.Accept()
		if err != nil {
			// select {
			// case <-s.getDoneChan():
			// 	return ErrServerClosed
			// default:
			// }
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if tempDelay > time.Second {
					tempDelay = time.Second
				}
				time.Sleep(tempDelay)
				continue
			}
			return err
		}

		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {

}
