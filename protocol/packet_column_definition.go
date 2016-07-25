package protocol

type PacketColumnDefinition struct {
	Packet       *Packet
	catalog      string
	schema       string
	table        string
	orgTable     string
	name         string
	orgName      string
	characterSet byte
	columnLength uint32
	columnType   uint8
	status       uint16
	decimals     uint8
}

func (p *PacketColumnDefinition) ToPacket() error {
	data := make([]byte, 4, 512)
	data = append(data, p.catalog...)
	data = append(data, 0)
	data = append(data, p.schema...)
	data = append(data, 0)
	data = append(data, p.table...)
	data = append(data, 0)
	data = append(data, p.orgTable...)
	data = append(data, 0)
	data = append(data, p.name...)
	data = append(data, 0)
	data = append(data, p.orgName...)
	data = append(data, 0)
	data = append(data, 0x0C)
	data = append(data, uint8(p.characterSet))
	data = append(data, byte(p.columnLength), byte(p.columnLength>>8), byte(p.columnLength>>16), byte(p.columnLength>>24))
	data = append(data, p.columnType)
	data = append(data, byte(p.status), byte(p.status>>8))
	data = append(data, p.decimals)
	data = append(data, byte(0), byte(0))
	return p.Packet.writePacket(data)
}

func DefaultProxyColumnDefinition(packet *Packet) *PacketColumnDefinition {
	p := &PacketColumnDefinition{}
	p.Packet = packet
	p.catalog = "def"
	p.schema = ""
	p.table = ""
	p.orgTable = ""
	p.name = "@@version_commen"
	p.orgName = ""
	p.characterSet = DefaultCollationID
	p.columnLength = uint32(len(VersionComment)) ///MariaDB Server
	p.columnType = MYSQL_TYPE_VAR_STRING
	p.status = uint16(0)
	p.decimals = 0x1f
	return p
}
