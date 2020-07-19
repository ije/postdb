package postdb

import (
	"fmt"
	"regexp"
	"strings"
)

var nameReg = regexp.MustCompile(`^[a-z0-9\_\-\.]+$`)

type ClientDB struct {
	client *Client
	db     string
}

// DB opens a client database
func (c *Client) DB(name string) (*ClientDB, error) {
	name = strings.TrimSpace(strings.ToLower(name))
	if !nameReg.MatchString(name) {
		return nil, fmt.Errorf("invalid name '%s'", name)
	}
	return &ClientDB{db: name, client: c}, nil
}
