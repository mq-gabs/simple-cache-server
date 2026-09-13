package handler

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"scas/cache"
	"scas/handler/process"
	"scas/store"

	"libsscas/protocol"

	"go.uber.org/zap"
)

type Handler struct {
	conn      net.Conn
	reader    bufio.Reader
	newreader io.Reader
	store     *store.Store
	cache     *cache.Cache
	log       *zap.SugaredLogger
}

func New(conn net.Conn, store *store.Store, logger *zap.SugaredLogger) *Handler {
	if logger == nil {
		panic("logger must be provided")
	}

	return &Handler{
		conn:      conn,
		reader:    *bufio.NewReader(conn),
		newreader: conn,
		store:     store,
		log:       logger.Named("handler"),
	}
}

func (h *Handler) readHeader() (*protocol.Header, error) {
	buf := make([]byte, protocol.HeaderSize)

	n, err := io.ReadFull(h.conn, buf)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCannotReadHeader, err)
	}
	if n != protocol.HeaderSize {
		return nil, errors.Join(ErrCannotReadHeader, ErrNumberOfReadBytesIsDifferentThanExpected)
	}

	header, err := protocol.DecodeHeader(buf)
	if err != nil {
		return nil, errors.Join(ErrCannotReadHeader, err)
	}

	return header, nil
}

func (h *Handler) readPayload(size protocol.PayloadLength) ([]byte, error) {
	buf := make([]byte, size)

	n, err := io.ReadFull(h.conn, buf)
	if err != nil {
		return nil, errors.Join(ErrCannotReadPayload, err)
	}
	if uint32(n) != uint32(size) {
		return nil, errors.Join(ErrCannotReadPayload, ErrNumberOfReadBytesIsDifferentThanExpected)
	}

	return buf, nil
}

func (h *Handler) Handle(ctx context.Context) {
	defer h.conn.Close()

	for {
		header, err := h.readHeader()
		if err != nil {
			h.log.Error("failed to read header", zap.Error(err))
			return
		}

		payload, err := h.readPayload(header.PayloadLength)
		if err != nil {
			h.log.Error("failed to read payload", zap.Error(err))
			return
		}

		resp, err := process.Process(h.cache, header, payload)
		if err != nil {
			h.log.Error("failed to process data", zap.Error(err))
			continue
		}
		if len(resp) == 0 {
			continue
		}
		if err := h.write(resp); err != nil {
			h.log.Error("failed to write response", zap.Error(err))
		}
	}
}

func (h *Handler) write(content []byte) error {
	_, err := h.conn.Write(content)
	return err
}
