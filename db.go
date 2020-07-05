package postdb

import (
	bolt "go.etcd.io/bbolt"
)

// A DB to store posts
type DB struct {
	*bolt.DB
}
