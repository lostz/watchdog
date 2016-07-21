package protocol

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
	authPluginData  string
	characterSet    byte
	status          uint16
	authPluginName  string
}

//ToPacket PacketHandshake struct to []byte
func (p PacketHandshake) ToPacket() error {
	data := make([]byte, 4, 128)
	data = append(data, 10)
	data = append(data, p.serverVersion...)
	data = append(data, p.protocolVersion)
	data = append(data, byte(p.connectionID), byte(p.connectionID>>8), byte(p.connectionID>>16), byte(p.connectionID>>24))
	data = append(data, p.authPluginData[0:8]...)
	data = append(data, 0x00)
	data = append(data, byte(p.Packet.capability), byte(p.Packet.capability>>8))
	data = append(data, uint8(p.characterSet))
	data = append(data, byte(p.status), byte(p.status>>8))
	data = append(data, byte(p.Packet.capability>>16), byte(p.Packet.capability>>24))
	data = append(data, 0x15)
	data = append(data, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0)
	data = append(data, p.authPluginData[8:]...)
	data = append(data, 0)
	return p.Packet.writePacket(data)
}

// NewPacketHandShake return default handShake
func NewPacketHandShake(connectionID uint32, salt string, packet *Packet) *PacketHandshake {
	p := &PacketHandshake{}
	p.Packet = packet
	p.connectionID = connectionID
	p.protocolVersion = MinProtocolVersion
	p.serverVersion = ServerVersion
	p.authPluginData = salt
	p.Packet.capability = DefaultCapability
	p.characterSet = DefaultCollationID
	p.status = SERVER_STATUS_AUTOCOMMIT
	p.authPluginName = AuthPluginName
	return p
}
