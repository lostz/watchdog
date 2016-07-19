package protocol

type PacketComStmtPrepareOK struct {
	Packet
	status       uint8
	statementID  uint32
	numColumns   uint16
	numParams    uint16
	warningCount uint16
}
