package process

import (
	"bytes"
	"libsscas/protocol"
	"scas/cache"
)

const (
	setSeparator         = 0x1e
	expectedLenAfterSpli = 2
	minKeyLength         = 4
	minValueLength       = 4
)

func processSet(c cache.Setter, flags protocol.Flag, payload []byte) ([]byte, error) {
	parts := bytes.Split(payload, []byte{setSeparator})
	if len(parts) != 2 {
		return nil, ErrInvalidSplittedLength
	}

	key := string(parts[0])
	if len(key) < minKeyLength {
		return nil, ErrInvalidKeyLength
	}

	value := parts[1]
	if len(value) < minValueLength {
		return nil, ErrInvalidValueLength
	}

	c.Set(key, value)

	return nil, nil
}
