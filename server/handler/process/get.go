package process

import (
	"libsscas/protocol"
	"scas/cache"
)

func processGet(c cache.Getter, flags protocol.Flag, payload []byte) ([]byte, error) {
	key := string(payload)

	value, ok := c.Get(key)
	if !ok {
		return nil, ErrKeyNotFound
	}

	return value, nil
}
