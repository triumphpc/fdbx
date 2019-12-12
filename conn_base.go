package fdbx

import (
	"encoding/binary"

	"github.com/apple/foundationdb/bindings/go/src/fdb"
)

func newBaseConn(db uint16) *baseConn {
	return &baseConn{
		db: db,
	}
}

type baseConn struct {
	db uint16
}

// ********************** Private **********************

func (c *baseConn) key(typeID uint16, parts ...[]byte) fdb.Key {
	mem := 4

	for i := range parts {
		mem += len(parts[i])
	}

	key := make(fdb.Key, 4, mem)

	binary.BigEndian.PutUint16(key[0:2], c.db)
	binary.BigEndian.PutUint16(key[2:4], typeID)

	for i := range parts {
		if len(parts[i]) > 0 {
			key = append(key, parts[i]...)
		}
	}

	return key
}

func (c *baseConn) rkey(rec Record) fdb.Key {
	rid := rec.FdbxID()
	rln := []byte{byte(len(rid))}
	return c.key(rec.FdbxType().ID, rid, rln)
}