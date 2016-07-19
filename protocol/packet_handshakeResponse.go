package protocol

// PacketHandshakeResponse 41
type PacketHandshakeResponse struct {
	sequenceID     uint8
	capability     uint32
	maxPacketSize  uint32
	characterSet   byte
	username       string
	authResponse   string
	database       string
	authPluginName string
	attributes     Attributes
}
