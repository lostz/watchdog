package protocol

type PacketResultsetRow struct {
	Packet *Packet
	data   []byte
}

func (p *PacketResultsetRow) SetData(data []byte) {
	p.data = data
}

func (p *PacketResultsetRow) ToPacket() error {
	data := make([]byte, 4, 512)
	if len(p.data) == 0 {
		data = append(data, 0xfb)
	} else {
		data = append(data, PutLengthEncodedString(p.data)...)
	}
	return p.Packet.writePacket(data)
}

func NewPacketResultsetRow(packet *Packet) *PacketResultsetRow {
	return &PacketResultsetRow{
		Packet: packet,
	}
}
