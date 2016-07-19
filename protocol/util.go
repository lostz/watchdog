package protocol

import (
	"encoding/hex"
	"io"
	"math/rand"
	"net"
	"strings"
	"time"
)

func LengthEncodedInt(b []byte) (num uint64, isNull bool, n int) {
	switch b[0] {

	// 251: NULL
	case 0xfb:
		n = 1
		isNull = true
		return

	// 252: value of following 2
	case 0xfc:
		num = uint64(b[1]) | uint64(b[2])<<8
		n = 3
		return

	// 253: value of following 3
	case 0xfd:
		num = uint64(b[1]) | uint64(b[2])<<8 | uint64(b[3])<<16
		n = 4
		return

	// 254: value of following 8
	case 0xfe:
		num = uint64(b[1]) | uint64(b[2])<<8 | uint64(b[3])<<16 |
			uint64(b[4])<<24 | uint64(b[5])<<32 | uint64(b[6])<<40 |
			uint64(b[7])<<48 | uint64(b[8])<<56
		n = 9
		return
	}

	// 0-250: value of first byte
	num = uint64(b[0])
	n = 1
	return
}

func RandomBuf(size int) []byte {
	buf := make([]byte, size)
	rand.Seed(time.Now().UTC().UnixNano())
	min, max := 30, 127
	for i := 0; i < size; i++ {
		buf[i] = byte(min + rand.Intn(max-min))
	}
	return buf
}

func LengthEnodedString(b []byte) ([]byte, bool, int, error) {
	// Get length
	num, isNull, n := LengthEncodedInt(b)
	if num < 1 {
		return nil, isNull, n, nil
	}

	n += int(num)

	// Check data length
	if len(b) >= n {
		return b[n-int(num) : n], false, n, nil
	}
	return nil, false, n, io.EOF
}

func PutLengthEncodedInt(n uint64) []byte {
	switch {
	case n <= 250:
		return []byte{byte(n)}

	case n <= 0xffff:
		return []byte{0xfc, byte(n), byte(n >> 8)}

	case n <= 0xffffff:
		return []byte{0xfd, byte(n), byte(n >> 8), byte(n >> 16)}

	case n <= 0xffffffffffffffff:
		return []byte{0xfe, byte(n), byte(n >> 8), byte(n >> 16), byte(n >> 24),
			byte(n >> 32), byte(n >> 40), byte(n >> 48), byte(n >> 56)}
	}
	return nil
}

func HasFlag(value uint64, flag uint64) bool {
	return (value & flag) == flag
}

func GetNulTerminatedStringSize(value string) uint64 {
	return uint64(len(value)) + 1
}

func WritePacket(c net.Conn, p Packet) error {
	data := p.CompressPacket()
	n, err := c.Write(data)
	if err != nil {
		return ErrBadConn
	} else if n != len(data) {
		return ErrBadConn
	} else {
		p.AddSequenceID()
		return nil
	}
}

func StringToPacket(value string) (data []byte) {
	lines := strings.Split(value, "\n")
	data = make([]byte, 0, 16*len(lines))
	var values []string

	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		if len(line) < 51 {
			values = strings.Split(line, " ")
		} else {
			values = strings.Split(line[:51], " ")
		}
		for _, val := range values {
			i, _ := hex.DecodeString(val)
			data = append(data, i...)
		}
	}

	return data
}
