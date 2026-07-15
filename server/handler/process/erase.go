package process

import (
	"libsscas/protocol"
	"scas/cache"
)

func processDelete(c cache.Deleter, flags protocol.Flag, payload []byte) ([]byte, error) {
	key := string(payload)
	c.Delete(key)

	return nil, nil
}
