package postdb

type Server struct {
	running bool
	usersDB *DB
	dbsDB   *DB
	DBDir   string
	Port    uint16
}

func (s *Server) CreateUser(name string, password string) error {
	return nil
}

func (s *Server) Serve() error {
	return nil
}
