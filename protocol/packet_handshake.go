package protocol

import (
	"bytes"
	"encoding/binary"
)

//DefaultCapability ...
var DefaultCapability uint32 = CLIENT_LONG_PASSWORD | CLIENT_LONG_FLAG |
	CLIENT_CONNECT_WITH_DB | CLIENT_PROTOCOL_41 |
	CLIENT_TRANSACTIONS | CLIENT_SECURE_CONNECTION

// PacketHandshake  V10
type PacketHandshake struct {
	Packet          *Packet
	protocolVersion byte
	serverVersion   string
	connectionID    uint32
	capability      uint32
	authPluginData  string
	characterSet    byte
	status          uint16
	authPluginName  string
}

//ToPacket PacketHandshake struct to []byte
func (p PacketHandshake) ToPacket() error {
	data := make([]byte, 4, 128)
	data = append(data, p.protocolVersion)
	data = append(data, p.serverVersion...)
	data = append(data, 0)
	data = append(data, byte(p.connectionID), byte(p.connectionID>>8), byte(p.connectionID>>16), byte(p.connectionID>>24))
	data = append(data, p.authPluginData[0:8]...)
	data = append(data, 0)
	data = append(data, byte(p.capability), byte(p.capability>>8))
	data = append(data, uint8(p.characterSet))
	data = append(data, byte(p.status), byte(p.status>>8))
	data = append(data, byte(p.capability>>16), byte(p.capability>>24))
	data = append(data, 0x15)
	data = append(data, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0)
	data = append(data, p.authPluginData[8:]...)
	data = append(data, 0)
	//if p.Packet.capability&CLIENT_PLUGIN_AUTH > 0 {
	//	data = append(data, p.authPluginName...)
	//	data = append(data, 0)
	//}
	return p.Packet.writePacket(data)
}

func (p *PacketHandshake) FromPacket(data []byte) {
	pos := 0
	p.protocolVersion = data[pos]
	pos++
	p.serverVersion = string(data[pos : pos+bytes.IndexByte(data[pos:], 0x00)])
	pos += len(p.serverVersion) + 1
	p.connectionID = binary.LittleEndian.Uint32(data[pos : pos+4])
	pos += 4
	auth1 := data[pos : pos+8]
	pos += 8
	// (filler) always 0x00 [1 byte]
	pos++
	p.capability = uint32(binary.LittleEndian.Uint16(data[pos:pos+2]))<<8 | p.capability
	pos += 2
	p.characterSet = data[pos]
	pos++
	p.status = binary.LittleEndian.Uint16(data[pos : pos+2])
	pos += 2
	p.capability = uint32(binary.LittleEndian.Uint16(data[pos:pos+2]))<<16 | p.capability
	pos += 2
	authLen := data[pos]
	pos++
	pos += 10
	auth2 := data[pos : pos+int(authLen)-8]
	pos += int(authLen) - 8
	p.authPluginData = string(auth1) + string(auth2)
	if p.capability&CLIENT_PLUGIN_AUTH > 0 {
		p.authPluginName = string(data[pos : pos+bytes.IndexByte(data[pos:], 0x00)])
	}

}

// NewPacketHandShake return default handShake
func NewPacketHandShake(connectionID uint32, salt string, packet *Packet) *PacketHandshake {
	p := &PacketHandshake{}
	p.Packet = packet
	p.connectionID = connectionID
	p.protocolVersion = MinProtocolVersion
	p.serverVersion = ServerVersion
	p.authPluginData = salt
	p.capability = DefaultCapability
	p.characterSet = DefaultCollationID
	p.status = SERVER_STATUS_AUTOCOMMIT
	p.authPluginName = AuthPluginName
	return p
}
