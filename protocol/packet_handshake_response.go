package protocol

import (
	"bytes"
	"encoding/binary"
)

// PacketHandshakeResponse 41
type PacketHandshakeResponse struct {
	Packet         *Packet
	packetSize     int
	maxPacketSize  uint32
	capability     uint32
	characterSet   byte
	username       string
	authResponse   string
	database       string
	authPluginName string
	attributes     map[string]string
}

func (p *PacketHandshakeResponse) Username() string {
	return p.username
}

func (p *PacketHandshakeResponse) AuthResponse() string {
	return p.authResponse
}

func (p *PacketHandshakeResponse) FromPacket() error {
	data, err := p.Packet.readPacket()
	if err != nil {
		return err
	}
	if data[0] == ErrHeader {
		return ErrBadConn
	}
	length := len(data)
	p.attributes = make(map[string]string)
	pos := 0
	p.capability = binary.LittleEndian.Uint32(data[:4])
	pos += 4
	p.maxPacketSize = binary.LittleEndian.Uint32(data[pos : pos+4])
	pos += 4
	p.characterSet = data[pos]
	pos++
	pos += 23
	p.username = string(data[pos : pos+bytes.IndexByte(data[pos:], 0)])
	pos += len(p.username) + 1
	authLen := int(data[pos])
	pos++
	p.authResponse = string(data[pos : pos+authLen])
	pos += authLen
	if p.capability&CLIENT_CONNECT_WITH_DB > 0 {
		db := string(data[pos : pos+bytes.IndexByte(data[pos:], 0)])
		pos += len(db) + 1
		p.database = db
	}
	if pos < length {
		if p.capability&CLIENT_PLUGIN_AUTH > 0 {
			authPluginName := string(data[pos : pos+bytes.IndexByte(data[pos:], 0)])
			pos += len(authPluginName) + 1
			p.authPluginName = authPluginName
		}
	}
	if pos < length {
		if p.capability&CLIENT_CONNECT_ATTRS > 0 {
			keyValueLen := int(data[pos])
			pos++
			for keyValueLen > 0 {
				key, _, n, err := LengthEnodedString(data[pos:])
				if err != nil {
					return err
				}
				pos += n
				p.attributes[string(key)] = ""
				keyValueLen -= n
				value, _, n, err := LengthEnodedString(data[pos:])
				if err != nil {
					return err
				}
				pos += n
				p.attributes[string(key)] = string(value)
				keyValueLen -= n
			}
		}
	}
	return nil

}
