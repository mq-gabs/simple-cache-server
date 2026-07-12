package protocol

import "errors"

var (
	ErrInvalidHeaderSize = errors.New("invalid header size")
	ErrInvalidVersion    = errors.New("invalid version")
	ErrInvalidMagic      = errors.New("invalid magic")
	ErrInvalidCommand    = errors.New("invalid command")

	ErrCannotDecodeHeader = errors.New("cannot decode header")
)
