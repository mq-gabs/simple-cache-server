package process

import (
	"errors"
	"libsscas/protocol"
	"scas/cache"
)

func processFlush(c cache.Flusher, flags protocol.Flag, payload []byte) ([]byte, error) {
	if len(payload) > 0 {
		return nil, errors.Join(ErrCrash, ErrPayloadNotExpected)
	}

	c.Flush()
	return nil, nil
}
