package protocol

import "net"

type Packet struct {
	netConn    net.Conn
	buf        buffer
	sequenceID uint8
}

func NewPacket(netConn net.Conn) *Packet {
	return &Packet{
		buf:     newBuffer(netConn),
		netConn: netConn,
	}

}

func (p *Packet) Conn() net.Conn {
	return p.netConn
}

func (p *Packet) CleanSequenceId() {
	p.sequenceID = 0
}

func (p *Packet) readPacket() ([]byte, error) {
	var payload []byte
	for {
		// Read packet header
		data, err := p.buf.readNext(4)
		if err != nil {
			p.Close()
			return nil, ErrBadConn
		}

		// Packet Length [24 bit]
		pktLen := int(uint32(data[0]) | uint32(data[1])<<8 | uint32(data[2])<<16)

		if pktLen < 1 {
			p.Close()
			return nil, ErrBadConn
		}

		// Check Packet Sync [8 bit]
		if data[3] != p.sequenceID {
			if data[3] > p.sequenceID {
				return nil, ErrMalformPacket
			}
			return nil, ErrMalformPacket
		}
		p.sequenceID++

		// Read packet body [pktLen bytes]
		data, err = p.buf.readNext(pktLen)
		if err != nil {
			p.Close()
			return nil, ErrBadConn
		}

		isLastPacket := (pktLen < MaxPayloadLen)

		// Zero allocations for non-splitting packets
		if isLastPacket && payload == nil {
			return data, nil
		}

		payload = append(payload, data...)

		if isLastPacket {
			return payload, nil
		}
	}
}

func (p *Packet) WritePacket(data []byte) error {
	return p.writePacket(data)
}

func (p *Packet) writePacket(data []byte) error {
	pktLen := len(data) - 4
	for {
		var size int
		if pktLen >= MaxPayloadLen {
			data[0] = 0xff
			data[1] = 0xff
			data[2] = 0xff
			size = MaxPayloadLen
		} else {
			data[0] = byte(pktLen)
			data[1] = byte(pktLen >> 8)
			data[2] = byte(pktLen >> 16)
			size = pktLen
		}
		data[3] = p.sequenceID

		n, err := p.netConn.Write(data[:4+size])
		if err == nil && n == 4+size {
			p.sequenceID++
			if size != MaxPayloadLen {
				return nil
			}
			pktLen -= size
			data = data[size:]
			continue
		}

		return ErrBadConn
	}
}

func (p *Packet) Close() (err error) {
	// Makes Close idempotent
	if p.netConn != nil {
		q := &PacketComQuit{p}
		err := q.WritePacket()
		p.netConn.Close()
		return err
	}
	p.netConn = nil

	return
}
