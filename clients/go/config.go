package scas

import (
	"bufio"
	"fmt"
	"net"
)

type Config struct {
	host string
	port uint16
}

func CreateConnection(config *Config) (*Connection, error) {
	if config == nil {
		config = &Config{}
	}
	if config.host == "" {
		config.host = "127.0.0.1"
	}
	if config.port == 0 {
		config.port = 9012
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%v", config.host, config.port))

	if err != nil {
		return nil, fmt.Errorf("cannot create connection: %v", err)
	}

	return &Connection{
		conn:   conn,
		reader: *bufio.NewReader(conn),
	}, nil
}

