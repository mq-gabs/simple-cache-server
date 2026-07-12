package protocol

import (
	"encoding/binary"
	"errors"
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

	if err := isValidMagic(Magic(binary.BigEndian.Uint16(buf[0:MagicSize]))); err != nil {
		return nil, errors.Join(ErrCannotDecodeHeader, err)
	}

	v := Version(buf[offsetVersion])
	if err := isValidVersion(v); err != nil {
		return nil, errors.Join(ErrCannotDecodeHeader, err)
	}

	c := Command(buf[offsetCommand])
	if err := isValidCommand(c); err != nil {
		return nil, errors.Join(ErrCannotDecodeHeader, err)
	}

	f := Flag(buf[offsetFlags])

	pLen := PayloadLength(binary.BigEndian.Uint32((buf[offsetPayloadLength : offsetPayloadLength+PayloadLengthSize])))

	h := &Header{
		Version:       v,
		Command:       c,
		Flags:         f,
		PayloadLength: pLen,
	}

	return h, nil
}
