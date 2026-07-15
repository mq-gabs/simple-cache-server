package process

import (
	"errors"
	"libsscas/protocol"
	"libsscas/protocol/validate"
	"scas/cache"
)

func processDelete(c cache.Deleter, flags protocol.Flag, payload []byte) ([]byte, error) {
	key := string(payload)

	if err := validate.IsValidKey(key); err != nil {
		return nil, errors.Join(ErrCrash, err)
	}

	c.Delete(key)

	return nil, nil
}
