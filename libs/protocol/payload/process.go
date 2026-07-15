package payload

import (
	"bytes"
)

const (
	payloadSetSeparator = 0x1e
)

func SplitPayloadSet(payload []byte) ([]byte, []byte, error) {
	parts := bytes.Split(payload, []byte{payloadSetSeparator})

	if len(parts) != 2 {
		return nil, nil, ErrInvalidPayload
	}

	key := parts[0]
	value := parts[1]

	return key, value, nil
}
