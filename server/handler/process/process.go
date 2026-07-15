package process

import (
	"fmt"
	"libsscas/protocol"
	"scas/cache"
)

func Process(c *cache.Cache, h *protocol.Header, payload []byte) (resp []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%w: %v", ErrProcessPanic, r)
			resp = nil
		}
	}()

	switch h.Command {
	case protocol.CmdGet:
		return processGet(c, h.Flags, payload)
	case protocol.CmdSet:
		return processSet(c, h.Flags, payload)
	case protocol.CmdDelete:
		return processDelete(c, h.Flags, payload)
	case protocol.CmdFlush:
		return processFlush(c, h.Flags, payload)
	default:
		return nil, ErrCommandNotImplemented
	}
}
