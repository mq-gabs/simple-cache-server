package process

import (
	"errors"
	"libsscas/protocol"
	"libsscas/protocol/validate"
	"scas/cache"
)

func processGet(c cache.Getter, flags protocol.Flag, payload []byte) ([]byte, error) {
	key := string(payload)

	if err := validate.IsValidKey(key); err != nil {
		return nil, errors.Join(ErrCrash, err)
	}

	value, ok := c.Get(key)
	if !ok {
		return nil, ErrKeyNotFound
	}

	return value, nil
}
