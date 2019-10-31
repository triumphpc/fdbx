package fdbx

// Supported FoundationDB client versions
const (
	ConnVersion610 = 610
)

// TxHandler -
type TxHandler func(DB) error

// Conn - database connection (as stored database index)
type Conn interface {
	Key(ctype uint16, id []byte) ([]byte, error)
	MKey(Model) ([]byte, error)

	Tx(TxHandler) error
}

// NewConn - makes new connection with specified client version
func NewConn(db, version uint16) (Conn, error) {
	// default 6.1.х
	if version == 0 {
		version = ConnVersion610
	}

	switch version {
	case ConnVersion610:
		return newV610Conn(db)
	}

	return nil, ErrUnknownVersion
}