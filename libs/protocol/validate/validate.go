package validate

const (
	minKeyLength   = 4
	minValueLength = 4
)

func IsValidKey(key string) error {
	if len(key) < minKeyLength {
		return ErrInvalidKeyLength
	}

	return nil
}

func IsValidValue(value []byte) error {
	if len(value) < minValueLength {
		return ErrInvalidValueLength
	}

	return nil
}
