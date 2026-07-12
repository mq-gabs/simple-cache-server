package protocol

import "fmt"

type Magic uint16

const (
	MagicSize       = 2
	SCAS      Magic = 0x5ca5
)

func isValidMagic(m Magic) error {
	switch m {
	case SCAS:
		return nil
	default:
		return fmt.Errorf("%w: %x", ErrInvalidMagic, m)
	}
}
