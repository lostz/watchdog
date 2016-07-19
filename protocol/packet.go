package protocol

import (
	"bytes"
	"compress/flate"
)

type Packet interface {
	SequenceID() uint8
	CompressPacket() []byte
	AddSequenceID()
	CleanSequenceID()
}

func CompressPacket(sequenceId uint8, input []byte) (output []byte) {
	var compressedPayloadLength int
	var compressedPayload []byte
	var uncompressedPayloadLength int
	uncompressedPayloadLength = len(input)
	if uncompressedPayloadLength < MIN_COMPRESS_LENGTH {
		compressedPayloadLength = uncompressedPayloadLength
		uncompressedPayloadLength = 0
		compressedPayload = input
	} else {
		var b bytes.Buffer
		w, _ := flate.NewWriter(&b, 9)
		w.Write(input)
		w.Close()
		compressedPayload = b.Bytes()
		compressedPayloadLength = len(compressedPayload)
	}
	output = make([]byte, 0, compressedPayloadLength+7)
	output = append(output, byte(compressedPayloadLength))
	output = append(output, byte(compressedPayloadLength>>8))
	output = append(output, byte(compressedPayloadLength>>16))
	output = append(output, sequenceId)
	output = append(output, byte(uncompressedPayloadLength))
	output = append(output, byte(uncompressedPayloadLength>>8))
	output = append(output, byte(uncompressedPayloadLength>>16))
	output = append(output, compressedPayload...)
	return output
}
