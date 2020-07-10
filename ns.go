package postdb

import (
	"encoding/base64"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

// A NSDB to store posts with namespaces
type NSDB struct {
	*DB
	lock       sync.RWMutex
	dbpath     string
	namespaces map[string]*DB
}

func New(path string) (ns *NSDB, err error) {
	dbpath, err := filepath.Abs(path)
	if err != nil {
		return
	}

	err = os.MkdirAll(dbpath, 0755)
	if err != nil {
		return
	}

	db, err := Open(strings.TrimSuffix(dbpath, "/")+"/public.db", 0666)
	if err != nil {
		return
	}

	ns = &NSDB{
		DB:         db,
		dbpath:     dbpath,
		namespaces: map[string]*DB{},
	}
	return
}

func (ns *NSDB) Namespace(name string) (db *DB, err error) {
	if name == "public" {
		return ns.DB, nil
	}

	ns.lock.RLock()
	db, ok := ns.namespaces[name]
	ns.lock.RUnlock()
	if !ok {
		db, err = Open(path.Join(ns.dbpath, base64.URLEncoding.EncodeToString([]byte(name))+".db"), 0666)
		if err != nil {
			return
		}
		ns.lock.Lock()
		ns.namespaces[name] = db
		ns.lock.Unlock()
	}
	return
}

func (ns *NSDB) Close() (err error) {
	ns.lock.Lock()
	defer ns.lock.Unlock()

	for name, db := range ns.namespaces {
		delete(ns.namespaces, name)
		db.Close()
	}
	return ns.DB.Close()
}
