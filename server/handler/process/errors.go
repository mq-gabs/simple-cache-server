package process

import "errors"

var (
	ErrCrash = errors.New("crash")

	ErrProcessPanic          = errors.New("process panic")
	ErrCommandNotImplemented = errors.New("command not implemented")
	ErrKeyNotFound           = errors.New("key not found")
	ErrInvalidSplittedLength = errors.New("invalid splitted length")
	ErrPayloadNotExpected    = errors.New("payload not expected")
)
