package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
)

const (
	bSep byte = byte(10)

	bReqGet byte = byte(48)	
	bReqSet byte = byte(49)
	bReqErase byte = byte(50)

	bRespErrorEmpty byte = byte(48)
	bRespErrorNotEmpty byte = byte(49)
	bRespSuccessEmpty byte = byte(50)
	bRespSuccessNotEmpty byte = byte(51)
)

type Handler struct {
	conn   net.Conn
	reader bufio.Reader
	store *Store
}

func NewHandler(conn net.Conn, store *Store) *Handler {
	return &Handler{
		conn: conn,
		reader: *bufio.NewReader(conn),
		store: store,
	}
}

func (h *Handler) readNext() ([]byte, error) {
	return h.reader.ReadBytes(bSep)
}

func (h *Handler) Handle() {
	defer h.conn.Close()

	for  {
		actHead, err := h.readNext()
		
		if err != nil {
			msg := fmt.Sprintf("Cannot read action: %v", err)
			log.Print(msg)
			h.respError(errors.New(msg))
			return
		}

		if len(actHead) > 2 || len(actHead) < 2 {
			msg := fmt.Sprintf("Action is invalid: %s", actHead)
			log.Print(msg)
			h.respError(errors.New(msg))
			return
		}

		switch actHead[0] {
			case bReqGet:
				h.get()		
			case bReqSet:
				h.set()
			case bReqErase:
				h.erase()
			default:
				log.Printf("Action not mapped: %v", actHead)
		}
	}
}

func (h *Handler) get() {
	key, err := h.readNext()

	if err != nil {
		msg := fmt.Sprintf("Cannot read key: %v", err)
		log.Print(msg)
		h.respError(errors.New(msg))
		return
	}

	sKey := string(key)

	value, err := h.store.Get(sKey)

	if err != nil {
		log.Printf("Error while getting key: %v", err)
		h.respError(err)
		return
	}
	
	log.Printf("Value: %s", value)
	h.respSuccess(value)
}

func (h *Handler) set() {
	key, err := h.readNext()

	if err != nil {
		msg := fmt.Sprintf("Cannot read key: %v", err)
		log.Print(msg)
		h.respError(errors.New(msg))
		return
	}
	
	sKey := string(key)

	value, err := h.readNext()

	if err != nil {
		log.Printf("Cannot read : %v", err)
	}

	err = h.store.Set(sKey, value)

	if err != nil {
		log.Printf("Error while saving: %v", err)
		h.respError(err)
		return
	}

	log.Println("Value saved!")
	h.respSuccess(nil)
}

func (h *Handler) erase() {
	key, err := h.readNext()

	if err != nil {
		msg := fmt.Sprintf("Cannot read key: %v", err)
		log.Print(msg)
		h.respError(errors.New(msg))
		return
	}
	
	sKey := string(key)

	err = h.store.Erase(sKey)
	if err != nil {
		log.Printf("Error while erasing: %v", err)
		h.respError(err)
		return
	}

	log.Println("Value erased!")
	h.respSuccess(nil)
}

func (h *Handler) respSuccess(content []byte) {
	if content == nil {
		h.conn.Write([]byte{bRespSuccessEmpty, bSep})
		return
	}

	bContent := JoinByte(content, bSep)

	h.conn.Write(JoinBytes([]byte{bRespSuccessNotEmpty, bSep}, bContent))
}

func (h *Handler) respError(err error) {
	if err == nil {
		h.conn.Write([]byte{bRespErrorEmpty, bSep})
		return
	}

	errBytes := JoinByte([]byte(err.Error()), bSep) 

	h.conn.Write(JoinBytes([]byte{bRespErrorNotEmpty, bSep}, errBytes))
}
