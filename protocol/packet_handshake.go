package protocol

import "math"

//DefaultCapability ...
var DefaultCapability uint32 = CLIENT_LONG_PASSWORD | CLIENT_LONG_FLAG |
	CLIENT_CONNECT_WITH_DB | CLIENT_PROTOCOL_41 |
	CLIENT_TRANSACTIONS | CLIENT_SECURE_CONNECTION

// PacketHandshake  V10
type PacketHandshake struct {
	sequenceID      uint8
	protocolVersion byte
	serverVersion   string
	connectionID    uint32
	authPluginData  string
	capability      uint32
	characterSet    byte
	status          uint16
	authPluginName  string
}

func (p PacketHandshake) SequenceID() uint8 {
	return p.sequenceID
}

func (p PacketHandshake) GetPacketSize() (size uint64) {
	size++
	size += GetNulTerminatedStringSize(p.serverVersion)
	size += 4
	size += 8
	size++
	size += 2
	size++
	size += 2
	size += 2
	size++
	size += 10
	if HasFlag(uint64(p.capability), uint64(CLIENT_SECURE_CONNECTION)) {
		size += uint64(math.Max(13, float64(len(p.authPluginData)-8)))
	}
	if HasFlag(uint64(p.capability), uint64(CLIENT_PLUGIN_AUTH)) {
		size += GetNulTerminatedStringSize(p.authPluginName)
	}
	return size
}

func (p PacketHandshake) ToPacket() (data []byte) {
	size := p.GetPacketSize()
	data = make([]byte, 0, size+4)
	data = append(data, PutLengthEncodedInt(uint64(size))...)
	data = append(data, p.sequenceID)
	data = append(data, p.protocolVersion)
	data = append(data, p.serverVersion...)
	data = append(data, byte(p.connectionID), byte(p.connectionID>>8), byte(p.connectionID>>16), byte(p.connectionID>>24))
	data = append(data, p.authPluginData[0:8]...)
	data = append(data, 0x00)
	data = append(data, byte(p.capability), byte(p.capability>>8))
	data = append(data, uint8(p.characterSet))
	data = append(data, byte(p.status), byte(p.status>>8))
	data = append(data, byte(p.capability>>16), byte(p.capability>>24))
	data = append(data, 0x15)
	data = append(data, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0)
	data = append(data, p.authPluginData[8:]...)
	data = append(data, 0)
	return data
}

func (p PacketHandshake) CompressPacket() []byte {
	data := p.ToPacket()
	return CompressPacket(p.sequenceID, data)
}

func (p PacketHandshake) AddSequenceID() {
	p.sequenceID++
}

func (p PacketHandshake) CleanSequenceID() {
	p.sequenceID = 0
}

// NewPacketHandShake return default handShake
func NewPacketHandShake(connectionID uint32, salt string) *PacketHandshake {
	p := &PacketHandshake{}
	p.sequenceID = 0
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
