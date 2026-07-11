package protocol

import "errors"

var (
	ErrInvalidHeaderSize = errors.New("invalid header size")
	ErrInvalidMagicBytes = errors.New("invalid magic bytes")
)
