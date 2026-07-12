package process

import (
	"libsscas/protocol"
	"scas/cache"
)

func Process(c *cache.Cache, h *protocol.Header, payload []byte) ([]byte, error) {
	switch h.Command {
	case protocol.CmdGet:
		return processGet(c, h.Flags, payload, int(h.PayloadLength))
	case protocol.CmdSet:
		return processSet(c, h.Flags, payload, int(h.PayloadLength))
	case protocol.CmdDelete:
		return processDelete(c, h.Flags, payload, int(h.PayloadLength))
	case protocol.CmdFlush:
		return processFlush(c, h.Flags, payload, int(h.PayloadLength))
	default:
		return nil, ErrCommandNotImplemented
	}
}
