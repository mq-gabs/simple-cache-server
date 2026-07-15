package process

import "errors"

var (
	ErrProcessPanic          = errors.New("process panic")
	ErrCommandNotImplemented = errors.New("command not implemented")
	ErrKeyNotFound           = errors.New("key not found")
	ErrInvalidSplittedLength = errors.New("invalid splitted length")
	ErrInvalidKeyLength      = errors.New("invalid key length")
	ErrInvalidValueLength    = errors.New("invalid value length")
)
