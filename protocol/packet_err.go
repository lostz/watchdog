package protocol

import "fmt"

//PacketErr packet error
type PacketErr struct {
	Packet       *Packet
	errCode      uint16
	sqlState     string
	errorMessage string
}

func NewDefaultPacketErr(packet *Packet, errCode uint16, args ...interface{}) *PacketErr {
	p := &PacketErr{
		Packet:  packet,
		errCode: errCode,
	}
	if format, ok := MySQLErrName[errCode]; ok {
		p.errorMessage = fmt.Sprintf(format, args...)
	} else {
		p.errorMessage = fmt.Sprint(args...)
	}
	return p
}

func NewPacketErr(packet *Packet, errMessage string) *PacketErr {
	return &PacketErr{
		Packet:       packet,
		errCode:      0,
		errorMessage: errMessage,
	}
}

//ToPacket  *PacketErr to []bye
func (p *PacketErr) ToPacket() error {
	if p.errCode < 1000 {
		p.errCode = ER_UNKNOWN_ERROR
	}
	if s, ok := MySQLState[p.errCode]; ok {
		p.sqlState = s
	} else {
		p.sqlState = DEFAULT_MYSQL_STATE
	}
	data := make([]byte, 4, 16+len(p.errorMessage))
	data = append(data, ErrHeader)
	data = append(data, byte(p.errCode), byte(p.errCode>>8))
	data = append(data, '#')
	data = append(data, p.sqlState...)
	data = append(data, p.errorMessage...)
	return p.Packet.writePacket(data)
}
