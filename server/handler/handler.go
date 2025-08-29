package handler

import (
	"bufio"
	"encoding/binary"
	"net"
	"scas/store"
	"scas/utils"
)

const (
	bSep byte = byte(10)

	bReqGet        byte = byte(48)
	bReqSet        byte = byte(49)
	bReqErase      byte = byte(50)
	bReqSetWithTTL byte = byte(51)

	bRespErrorEmpty     byte = byte(48)
	bRespErrorContent   byte = byte(49)
	bRespSuccessEmpty   byte = byte(50)
	bRespSuccessContent byte = byte(51)
)

type Handler struct {
	conn   net.Conn
	reader bufio.Reader
	store  *store.Store
}

func New(conn net.Conn, store *store.Store) *Handler {
	return &Handler{
		conn:   conn,
		reader: *bufio.NewReader(conn),
		store:  store,
	}
}

func (h *Handler) readNext() ([]byte, error) {
	return h.reader.ReadBytes(bSep)
}

func (h *Handler) Handle() {
	defer h.conn.Close()

	for {
		actHead, err := h.readNext()

		if err != nil {
			h.respError(utils.FmtErr(errCannotReadAction, err))
			return
		}

		if len(actHead) != 2 {
			h.respError(utils.FmtErr(errActionIsInvalid, actHead))
			return
		}

		switch actHead[0] {
		case bReqGet:
			h.get()
		case bReqSet:
			h.set()
		case bReqErase:
			h.erase()
		case bReqSetWithTTL:
			h.setWithTTL()
		default:
			h.respError(utils.FmtErr(errActionDoesNotExists, actHead))
		}
	}
}

func (h *Handler) write(content []byte) error {
	_, err := h.conn.Write(content)
	return err
}

func (h *Handler) bytesJoin(bytes ...[]byte) []byte {
	res := utils.JoinBytes(bSep, bytes...)
	return append(res, bSep)
}

func (h *Handler) setWithTTL() {
	key, err := h.readNext()

	if err != nil {
		h.respError(utils.FmtErr(errCannotReadKey, err))
		return
	}

	keyStr := string(key)

	value, err := h.readNext()
	if err != nil {
		h.respError(utils.FmtErr(errCannotReadValue, err))
		return
	}

	ttl, err := h.readNext()
	if err != nil {
		h.respError(utils.FmtErr(errCannotReadTTL, err))
	}

	intTtl := binary.BigEndian.Uint32(ttl)

	if err := h.store.SetWithTTL(keyStr, value, intTtl); err != nil {
		h.respError(err)
	}

	h.respSuccess(nil)
}

func (h *Handler) get() {
	key, err := h.readNext()
	if err != nil {
		h.respError(utils.FmtErr(errCannotReadKey, err))
		return
	}

	keyStr := string(key)

	value, err := h.store.Get(keyStr)
	if err != nil {
		h.respError(utils.FmtErr(errCannotGetKey, err))
		return
	}

	h.respSuccess(*value)
}

func (h *Handler) set() {
	key, err := h.readNext()
	if err != nil {
		h.respError(utils.FmtErr(errCannotReadKey, err))
		return
	}

	keyStr := string(key)

	value, err := h.readNext()
	if err != nil {
		h.respError(utils.FmtErr(errCannotReadValue, err))
		return
	}

	if err = h.store.Set(keyStr, value); err != nil {
		h.respError(utils.FmtErr(errCannotSet, err))
		return
	}

	h.respSuccess(nil)
}

func (h *Handler) erase() {
	key, err := h.readNext()
	if err != nil {
		h.respError(utils.FmtErr(errCannotReadKey, err))
		return
	}

	keyStr := string(key)

	if err = h.store.Erase(keyStr); err != nil {
		h.respError(utils.FmtErr(errCannotErase, err))
		return
	}

	h.respSuccess(nil)
}

func (h *Handler) respSuccess(content []byte) {
	if content == nil {
		h.write([]byte{bRespSuccessEmpty, bSep})
		return
	}

	head := []byte{bRespSuccessContent}

	h.write(h.bytesJoin(head, content))
}

func (h *Handler) respError(err error) {
	if err == nil {
		h.write([]byte{bRespErrorEmpty, bSep})
		return
	}

	head := []byte{bRespErrorContent}
	errBytes := []byte(err.Error())

	h.write(h.bytesJoin(head, errBytes))
}
