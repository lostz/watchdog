package protocol

type PacketEOF struct {
	Packet       *Packet
	warningCount uint16
	status       uint16
}

func (p *PacketEOF) ToPacket() error {
	data := make([]byte, 4, 8)
	data = append(data, EOFHeader)
	data = append(data, byte(p.warningCount), byte(p.warningCount>>8))
	data = append(data, byte(p.status), byte(p.status>>8))
	return p.Packet.writePacket(data)
}

func (p *PacketEOF) SetStatus(status uint16) {
	p.status = status
}

func NewPacketEOF(packet *Packet) *PacketEOF {
	return &PacketEOF{
		Packet: packet,
	}

}
