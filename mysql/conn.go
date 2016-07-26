package mysql

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/lostz/watchdog/protocol"
)

type Conn struct {
	conn       net.Conn
	packet     *protocol.Packet
	capability uint32
	addr       string
	user       string
	password   string
	db         string
	collation  uint8
	charset    string
	salt       []byte
	status     uint16
}

func (c *Conn) Connect(addr string, user string, password string, db string) error {
	c.addr = addr
	c.user = user
	c.password = password
	c.db = db

	c.collation = uint8(33)
	c.charset = "utf8"
	return c.ReConnect()
}

func (c *Conn) ReConnect() error {
	if c.conn != nil {
		c.conn.Close()
	}

	n := "tcp"
	if strings.Contains(c.addr, "/") {
		n = "unix"
	}

	netConn, err := net.Dial(n, c.addr)
	if err != nil {
		return err
	}

	tcpConn := netConn.(*net.TCPConn)

	tcpConn.SetNoDelay(false)
	tcpConn.SetKeepAlive(true)
	c.conn = tcpConn
	c.packet = protocol.NewPacket(tcpConn)

	if err := c.readInitialHandshake(); err != nil {
		c.conn.Close()
		return err
	}
	if err := c.writeAuthHandshake(); err != nil {
		c.conn.Close()

		return err
	}
	if err := c.readOK(); err != nil {
		c.conn.Close()

		return err
	}
	return nil
}

func (c *Conn) ReadPacket() ([]byte, error) {
	d, err := c.packet.ReadPacket()
	return d, err
}

func (c *Conn) WritePacket(data []byte) error {
	err := c.packet.WritePacket(data)
	return err
}

func (c *Conn) readInitialHandshake() error {
	data, err := c.ReadPacket()
	if err != nil {
		return err
	}

	if data[0] == protocol.ErrHeader {
		return errors.New("read initial handshake error")
	}

	if data[0] < protocol.MinProtocolVersion {
		return fmt.Errorf("invalid protocol version %d, must >= 10", data[0])
	}

	//skip mysql version and connection id
	//mysql version end with 0x00
	//connection id length is 4
	pos := 1 + bytes.IndexByte(data[1:], 0x00) + 1 + 4

	c.salt = append(c.salt, data[pos:pos+8]...)

	//skip filter
	pos += 8 + 1

	//capability lower 2 bytes
	c.capability = uint32(binary.LittleEndian.Uint16(data[pos : pos+2]))

	pos += 2

	if len(data) > pos {
		//skip server charset
		//c.charset = data[pos]
		pos += 1

		c.status = binary.LittleEndian.Uint16(data[pos : pos+2])
		pos += 2

		c.capability = uint32(binary.LittleEndian.Uint16(data[pos:pos+2]))<<16 | c.capability

		pos += 2

		//skip auth data len or [00]
		//skip reserved (all [00])
		pos += 10 + 1

		// The documentation is ambiguous about the length.
		// The official Python library uses the fixed length 12
		// mysql-proxy also use 12
		// which is not documented but seems to work.
		c.salt = append(c.salt, data[pos:pos+12]...)
	}

	return nil
}

func (c *Conn) writeAuthHandshake() error {
	// Adjust client capability flags based on server support
	capability := protocol.CLIENT_PROTOCOL_41 | protocol.CLIENT_SECURE_CONNECTION |
		protocol.CLIENT_LONG_PASSWORD | protocol.CLIENT_TRANSACTIONS | protocol.CLIENT_LONG_FLAG

	capability &= c.capability

	//packet length
	//capbility 4
	//max-packet size 4
	//charset 1
	//reserved all[0] 23
	length := 4 + 4 + 1 + 23

	//username
	length += len(c.user) + 1

	//we only support secure connection
	auth := protocol.CalcPassword(c.salt, []byte(c.password))

	length += 1 + len(auth)

	if len(c.db) > 0 {
		capability |= protocol.CLIENT_CONNECT_WITH_DB

		length += len(c.db) + 1
	}

	c.capability = capability

	data := make([]byte, length+4)

	//capability [32 bit]
	data[4] = byte(capability)
	data[5] = byte(capability >> 8)
	data[6] = byte(capability >> 16)
	data[7] = byte(capability >> 24)

	//MaxPacketSize [32 bit] (none)
	//data[8] = 0x00
	//data[9] = 0x00
	//data[10] = 0x00
	//data[11] = 0x00

	//Charset [1 byte]
	data[12] = byte(c.collation)

	//Filler [23 bytes] (all 0x00)
	pos := 13 + 23

	//User [null terminated string]
	if len(c.user) > 0 {
		pos += copy(data[pos:], c.user)
	}
	//data[pos] = 0x00
	pos++

	// auth [length encoded integer]
	data[pos] = byte(len(auth))
	pos += 1 + copy(data[pos+1:], auth)

	// db [null terminated string]
	if len(c.db) > 0 {
		pos += copy(data[pos:], c.db)
		//data[pos] = 0x00
	}

	return c.WritePacket(data)
}

func (c *Conn) readOK() error {
	data, err := c.ReadPacket()
	if err != nil {
		return err
	}

	if data[0] == protocol.OkHeader {
		return c.handleOKPacket(data)
	} else if data[0] == protocol.ErrHeader {
		return c.handleErrorPacket(data)
	} else {
		return errors.New("invalid ok packet")
	}
}

func (c *Conn) handleOKPacket(data []byte) error {
	var n int
	var pos int = 1

	_, _, n = protocol.LengthEncodedInt(data[pos:])
	pos += n
	_, _, n = protocol.LengthEncodedInt(data[pos:])
	pos += n

	if c.capability&protocol.CLIENT_PROTOCOL_41 > 0 {
		c.status = binary.LittleEndian.Uint16(data[pos:])
		pos += 2

		//todo:strict_mode, check warnings as error
		//Warnings := binary.LittleEndian.Uint16(data[pos:])
		//pos += 2
	} else if c.capability&protocol.CLIENT_TRANSACTIONS > 0 {
		c.status = binary.LittleEndian.Uint16(data[pos:])
		pos += 2
	}
	return nil
}

func (c *Conn) handleErrorPacket(data []byte) error {

	var pos int = 1
	var state string
	var message string

	code := binary.LittleEndian.Uint16(data[pos:])
	pos += 2

	if c.capability&protocol.CLIENT_PROTOCOL_41 > 0 {
		//skip '#'
		pos++
		state = string(data[pos : pos+5])
		pos += 5
	}

	message = string(data[pos:])
	return errors.New(fmt.Sprintf("code: %d state: %s meessage: %s", code, state, message))

}

func (c *Conn) Close() error {
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
		c.salt = nil
	}

	return nil
}
