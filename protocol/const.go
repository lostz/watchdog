package protocol

const (
	//MinProtocolVersion 10
	MinProtocolVersion byte = 10
	// MaxPayloadLen  8388608
	MaxPayloadLen int = 1<<24 - 1
	//TimeFormat ...
	TimeFormat string = "2006-01-02 15:04:05"
	//ServerVersion handShake ...
	ServerVersion string = "5.5.19-lostz-0.1"
	//AuthPluginName handShake ...
	AuthPluginName string = "mysql_native_password"
	//DefaultCollationID utf8
	DefaultCollationID uint8 = 33
)

const (
	SERVER_STATUS_IN_TRANS             uint16 = 0x0001
	SERVER_STATUS_AUTOCOMMIT           uint16 = 0x0002
	SERVER_MORE_RESULTS_EXISTS         uint16 = 0x0008
	SERVER_STATUS_NO_GOOD_INDEX_USED   uint16 = 0x0010
	SERVER_STATUS_NO_INDEX_USED        uint16 = 0x0020
	SERVER_STATUS_CURSOR_EXISTS        uint16 = 0x0040
	SERVER_STATUS_LAST_ROW_SEND        uint16 = 0x0080
	SERVER_STATUS_DB_DROPPED           uint16 = 0x0100
	SERVER_STATUS_NO_BACKSLASH_ESCAPED uint16 = 0x0200
	SERVER_STATUS_METADATA_CHANGED     uint16 = 0x0400
	SERVER_QUERY_WAS_SLOW              uint16 = 0x0800
	SERVER_PS_OUT_PARAMS               uint16 = 0x1000
)

const (
	CLIENT_LONG_PASSWORD uint32 = 1 << iota
	CLIENT_FOUND_ROWS
	CLIENT_LONG_FLAG
	CLIENT_CONNECT_WITH_DB
	CLIENT_NO_SCHEMA
	CLIENT_COMPRESS
	CLIENT_ODBC
	CLIENT_LOCAL_FILES
	CLIENT_IGNORE_SPACE
	CLIENT_PROTOCOL_41
	CLIENT_INTERACTIVE
	CLIENT_SSL
	CLIENT_IGNORE_SIGPIPE
	CLIENT_TRANSACTIONS
	CLIENT_RESERVED
	CLIENT_SECURE_CONNECTION
	CLIENT_MULTI_STATEMENTS
	CLIENT_MULTI_RESULTS
	CLIENT_PS_MULTI_RESULTS
	CLIENT_PLUGIN_AUTH
	CLIENT_CONNECT_ATTRS
	CLIENT_PLUGIN_AUTH_LENENC_CLIENT_DATA
)

const (
	OKHeader          byte = 0x00
	ErrHeader         byte = 0xff
	EofHeader         byte = 0xfe
	LocalInFileHeader byte = 0xfb
)

const (
	comQuit byte = iota + 1
	comInitDB
	comQuery
	comFieldList
	comCreateDB
	comDropDB
	comRefresh
	comShutdown
	comStatistics
	comProcessInfo
	comConnect
	comProcessKill
	comDebug
	comPing
	comTime
	comDelayedInsert
	comChangeUser
	comBinlogDump
	comTableDump
	comConnectOut
	comRegisterSlave
	comStmtPrepare
	comStmtExecute
	comStmtSendLongData
	comStmtClose
	comStmtReset
	comSetOption
	comStmtFetch
)
