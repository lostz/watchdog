package protocol

type Context struct {
	capability          uint64
	prepared_statements map[uint32]PacketComStmtPrepareOK
}
