package protocol

type PacketComQuit struct {
	*Packet
}

func (p *PacketComQuit) ToPacket() (data []byte) {
	p.Packet.sequenceID = 0
	data = make([]byte, 5)
	data[4] = comQuit
	return data
}

func (p *PacketComQuit) WritePacket() error {
	return p.Packet.writePacket(p.ToPacket())
}
