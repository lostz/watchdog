package protocol

// PacketOk  packet ok
type PacketOk struct {
	Packet             *Packet
	affectedRows       uint64
	lastInsertID       uint64
	status             uint16
	warnings           uint16
	info               string
	sessionStateChange string
}

//ToPacket PacketOk to []byte
func (p *PacketOk) ToPacket() error {
	data := make([]byte, 4, 32)
	data = append(data, OkHeader)
	data = append(data, PutLengthEncodedInt(p.affectedRows)...)
	data = append(data, PutLengthEncodedInt(p.lastInsertID)...)
	data = append(data, byte(p.status), byte(p.status>>8))
	data = append(data, 0, 0)
	return p.Packet.writePacket(data)

}

//NewPacketOk return *PacketOk
func NewPacketOk(packet *Packet) *PacketOk {
	return &PacketOk{
		Packet: packet,
	}
}
