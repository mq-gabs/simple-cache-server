package protocol

import "fmt"

type Command byte

const (
	CommandSize = 1

	CmdInvalid Command = 0x00

	CmdGet    Command = 0x01
	CmdSet    Command = 0x02
	CmdDelete Command = 0x03
	CmdFlush  Command = 0x04
)

func isValidCommand(c Command) error {
	switch c {
	case CmdGet, CmdSet, CmdDelete, CmdFlush:
		return nil
	default:
		return fmt.Errorf("%w: %x", ErrInvalidCommand, c)
	}
}
