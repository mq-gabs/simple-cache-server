package process

import (
	"libsscas/protocol"
	"scas/cache"
)

func processFlush(c cache.Flusher, flags protocol.Flag, payload []byte) ([]byte, error) {
	c.Flush()
	return nil, nil
}
