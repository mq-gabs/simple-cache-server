package protocol

import (
	"encoding/binary"
	"fmt"
)

const (
	HeaderSize = MagicSize + VersionSize + CommandSize + FlagsSize + PayloadLengthSize

	offsetVersion       = 2
	offsetCommand       = 3
	offsetFlags         = 4
	offsetPayloadLength = 5
)

type Header struct {
	Version       Version
	Command       Command
	Flags         Flag
	PayloadLength PayloadLength
}

func DecodeHeader(buf []byte) (*Header, error) {
	if len(buf) != HeaderSize {
		return nil, fmt.Errorf("%w: %v", ErrInvalidHeaderSize, len(buf))
	}

	if SCAS != Magic(binary.BigEndian.Uint16(buf[0:MagicSize])) {
		return nil, ErrInvalidMagicBytes
	}

	h := &Header{
		Version:       Version(buf[offsetVersion]),
		Command:       Command(buf[offsetCommand]),
		Flags:         Flag(buf[offsetFlags]),
		PayloadLength: PayloadLength(binary.BigEndian.Uint32((buf[offsetPayloadLength : offsetPayloadLength+PayloadLengthSize]))),
	}

	return h, nil
}
