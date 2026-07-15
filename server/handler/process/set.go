package process

import (
	"errors"
	"libsscas/protocol"
	protopayload "libsscas/protocol/payload"
	"libsscas/protocol/validate"
	"scas/cache"
)

const (
	setSeparator         = 0x1e
	expectedLenAfterSpli = 2
)

func processSet(c cache.Setter, flags protocol.Flag, payload []byte) ([]byte, error) {
	keyBytes, value, err := protopayload.SplitPayloadSet(payload)
	if err != nil {
		return nil, errors.Join(ErrCrash, err)
	}
	key := string(keyBytes)

	if err := validate.IsValidKey(key); err != nil {
		return nil, errors.Join(ErrCrash, err)
	}

	if err := validate.IsValidValue(value); err != nil {
		return nil, errors.Join(ErrCrash, err)
	}

	c.Set(key, value)

	return nil, nil
}
