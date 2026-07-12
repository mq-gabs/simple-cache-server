package protocol

import "fmt"

type Version byte

const (
	VersionSize = 1

	Version0 Version = 0x00
)

func isValidVersion(v Version) error {
	switch v {
	case Version0:
		return nil
	default:
		return fmt.Errorf("%w: %x", ErrInvalidVersion, v)
	}
}
