package scas

import (
	"bufio"
	"errors"
	"fmt"
	"net"
)

type Connection struct {
	conn   net.Conn
	reader bufio.Reader
}

const (
	bSep byte = byte(10)

	bReqGet   byte = byte(48)
	bReqSet   byte = byte(49)
	bReqErase byte = byte(50)

	bRespErrorEmpty      byte = byte(48)
	bRespErrorNotEmpty   byte = byte(49)
	bRespSuccessEmpty    byte = byte(50)
	bRespSuccessNotEmpty byte = byte(51)
)

func (c *Connection) readNext() ([]byte, error) {
	return c.reader.ReadBytes(bSep)
}

func (c *Connection) readResponse() (string, error) {
	s, err := c.readNext()

	if err != nil {
		return "", fmt.Errorf("cannot read response status: %s", err)
	}

	if len(s) > 2 || len(s) < 2 {
		return "", fmt.Errorf("status is invalid: %v", s)
	}

	status := s[0]

	if status == bRespErrorEmpty {
		return "", errors.New("some error ocurred")
	}

	if status == bRespErrorNotEmpty {
		content, err := c.readNext()

		if err != nil {
			return "", fmt.Errorf("cannot read error content: %v", err)
		}

		return "", fmt.Errorf("server Error: %s", content)
	}

	if status == bRespSuccessEmpty {
		return "", nil
	}

	if status == bRespSuccessNotEmpty {
		content, err := c.readNext()

		if err != nil {
			return "", fmt.Errorf("cannot read success content: %s", err)
		}

		return string(content), nil
	}

	return "", fmt.Errorf("status not mapped: %v", status)
}

func (c *Connection) Get(key string) (string, error) {
	bKey := joinByte([]byte(key), bSep)

	_, err := c.conn.Write(joinBytes([]byte{bReqGet, bSep}, bKey))

	if err != nil {
		return "", fmt.Errorf("cannot get: %v", err)
	}

	return c.readResponse()
}

func (c *Connection) Set(key, value string) error {
	bKey := joinByte([]byte(key), bSep)
	bValue := joinByte([]byte(value), bSep)

	_, err := c.conn.Write(join2Bytes([]byte{bReqSet, bSep}, bKey, bValue))

	if err != nil {
		return fmt.Errorf("cannot set: %v", err)
	}

	_, err = c.readResponse()

	return err
}

func (c *Connection) Erase(key string) error {
	bKey := joinByte([]byte(key), bSep)

	_, err := c.conn.Write(joinBytes([]byte{bReqErase, bSep}, bKey))

	if err != nil {
		return fmt.Errorf("cannot erase: %v", err)
	}

	_, err = c.readResponse()

	return err
}

func (c *Connection) Close() error {
	return c.conn.Close()
}
