package validate

import "errors"

var (
	ErrInvalidKeyLength   = errors.New("invalid key length")
	ErrInvalidValueLength = errors.New("invalid value length")
)
