package protocol

type PacketColumnCount struct {
	Packet      *Packet
	columnCount int
}

func (p *PacketColumnCount) SetColumnCount(count int) {
	p.columnCount = 0
}

func (p *PacketColumnCount) ToPacket() error {
	data := make([]byte, 4, 512)
	data = append(data, byte(p.columnCount))
	return p.Packet.writePacket(data)
}

func (p *PacketColumnCount) FromPacket() error {
	data, err := p.Packet.readPacket()
	if err != nil {
		return err
	}
	p.columnCount = int(data[0])
	return nil
}

func NewPacketColumnCount(packet *Packet) *PacketColumnCount {
	return &PacketColumnCount{
		Packet: packet,
	}
}
