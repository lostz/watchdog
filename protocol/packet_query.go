package protocol

//PacketQuery packet query
type PacketQuery struct {
	Packet *Packet
	query  string
}

//FromPacket []byte to PacketQuery
func (p *PacketQuery) FromPacket() error {
	data, err := p.Packet.readPacket()
	if err != nil {
		return err
	}
	p.query = string(data[1:])
	return nil
}

//Query return query
func (p *PacketQuery) Query() string {
	return p.query
}

//NewPacketQuery return *PacketQuery
func NewPacketQuery(packet *Packet) *PacketQuery {
	return &PacketQuery{
		Packet: packet,
	}
}
