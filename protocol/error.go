package protocol

import "errors"

var (
	ErrBadConn       = errors.New("connection was bad")
	ErrMalformPacket = errors.New("Malform packet error")
	ErrTxDone        = errors.New("sql: Transaction has already been committed or rolled back")
)
