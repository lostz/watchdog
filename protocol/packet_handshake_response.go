package protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// PacketHandshakeResponse 41
type PacketHandshakeResponse struct {
	sequenceID     uint8
	capability     uint32
	maxPacketSize  uint32
	characterSet   byte
	username       string
	authResponse   string
	database       string
	authPluginName string
}

func (p *PacketHandshakeResponse) FromPacket(data []byte) {
	length := int(uint32(data[0]) | uint32(data[1])<<8 | uint32(data[2])<<16)
	fmt.Println(length)
	p.sequenceID = uint8(data[3])
	data = data[4:]
	pos := 0
	p.capability = binary.LittleEndian.Uint32(data[:4])
	pos += 4
	p.maxPacketSize = binary.LittleEndian.Uint32(data[pos : pos+4])
	pos += 4
	p.characterSet = data[pos]
	pos++
	pos += 23
	p.username = string(data[pos : pos+bytes.IndexByte(data[pos:], 0)])
	pos += len(p.username) + 1
	authLen := int(data[pos])
	pos++
	p.authResponse = string(data[pos : pos+authLen])
	pos += authLen
	if p.capability&CLIENT_CONNECT_WITH_DB > 0 {
		db := string(data[pos : pos+bytes.IndexByte(data[pos:], 0)])
		pos += len(db) + 1
		p.database = db
	}
	fmt.Println(pos)

}
