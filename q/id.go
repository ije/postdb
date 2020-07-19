package q

import (
	"github.com/rs/xid"
)

// ObjectID represents a unique object id powered by xid(https://github.com/rs/xid)
type ObjectID [12]byte

var nilObjectID ObjectID

// NewID returns a new ID
func NewID() ObjectID {
	return ObjectID(xid.New())
}

// ID returns an ObjectID
func ID(id string) ObjectID {
	xid, err := xid.FromString(id)
	if err != nil {
		return nilObjectID
	}
	return ObjectID(xid)
}

// String returns a base32 hex lowercased with no padding representation of the id (char set is 0-9, a-v).
func (id ObjectID) String() string {
	return xid.ID(id).String()
}

// Bytes returns the byte array representation of `ID`
func (id ObjectID) Bytes() []byte {
	return id[:]
}

// IsNil returns true if this is a "nil" ID
func (id ObjectID) IsNil() bool {
	return id == nilObjectID
}

// QueryType implements the Query interface
func (id ObjectID) QueryType() string {
	return "id"
}
