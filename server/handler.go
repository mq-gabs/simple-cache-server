package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"net"
)

const (
	bSep byte = byte(10)

	bReqGet        byte = byte(48)
	bReqSet        byte = byte(49)
	bReqErase      byte = byte(50)
	bReqSetWithTTL byte = byte(51)

	bRespErrorEmpty      byte = byte(48)
	bRespErrorNotEmpty   byte = byte(49)
	bRespSuccessEmpty    byte = byte(50)
	bRespSuccessNotEmpty byte = byte(51)
)

type Handler struct {
	conn   net.Conn
	reader bufio.Reader
	store  *Store
}

func NewHandler(conn net.Conn, store *Store) *Handler {
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
			h.respError(fmt.Errorf("cannot read action: %v", err))
			return
		}

		if len(actHead) != 2 {
			h.respError(fmt.Errorf("action is invalid: %v", actHead))
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
			h.respError(fmt.Errorf("action not mapped: %v", actHead))
		}
	}
}

func (h *Handler) setWithTTL() {
	key, err := h.readNext()

	if err != nil {
		h.respError(fmt.Errorf("cannot read key: %v", err))
		return
	}

	sKey := string(key)

	value, err := h.readNext()

	if err != nil {
		h.respError(fmt.Errorf("cannot read value: %v", err))
		return
	}

	ttl, err := h.readNext()
	intTtl := binary.BigEndian.Uint32(ttl)

	if err != nil {
		h.respError(fmt.Errorf("cannot read time to live: %v", err))
	}

	if err := h.store.SetWithTTL(sKey, value, intTtl); err != nil {
		h.respError(err)
	}

	h.respSuccess(nil)
}

func (h *Handler) get() {
	key, err := h.readNext()

	if err != nil {
		h.respError(fmt.Errorf("cannot read key: %v", err))
		return
	}

	sKey := string(key)

	value, err := h.store.Get(sKey)

	if err != nil {
		h.respError(fmt.Errorf("error while getting key: %v", err))
		return
	}

	h.respSuccess(*value)
}

func (h *Handler) set() {
	key, err := h.readNext()

	if err != nil {
		h.respError(fmt.Errorf("cannot read key: %v", err))
		return
	}

	sKey := string(key)

	value, err := h.readNext()

	if err != nil {
		h.respError(fmt.Errorf("Cannot read value: %v", err))
		return
	}

	err = h.store.Set(sKey, value)

	if err != nil {
		h.respError(fmt.Errorf("error while setting: %v", err))
		return
	}

	h.respSuccess(nil)
}

func (h *Handler) erase() {
	key, err := h.readNext()

	if err != nil {
		h.respError(fmt.Errorf("cannot read key: %v", err))
		return
	}

	sKey := string(key)

	err = h.store.Erase(sKey)
	if err != nil {
		h.respError(fmt.Errorf("error while erasing: %v", err))
		return
	}

	h.respSuccess(nil)
}

func (h *Handler) respSuccess(content []byte) {
	if content == nil {
		h.conn.Write([]byte{bRespSuccessEmpty, bSep})
		return
	}

	// bContent := joinByte(content, bSep)

	head := []byte{bRespSuccessNotEmpty, bSep}

	h.conn.Write(joinBytes([][]byte{head, content}))
}

func (h *Handler) respError(err error) {
	if err == nil {
		h.conn.Write([]byte{bRespErrorEmpty, bSep})
		return
	}

	head := []byte{bRespErrorNotEmpty, bSep}
	errBytes := append([]byte(err.Error()), bSep)

	h.conn.Write(joinBytes([][]byte{head, errBytes}))
}
