package mysql

type BackendConn interface {
	ReadPacket() ([]byte, error)
	WritePacket(data []byte) error
	Close() error
}
