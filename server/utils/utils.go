package utils

import "bytes"

func JoinBytes(sep byte, slices ...[]byte) []byte {
	return bytes.Join(slices, []byte{sep})
}
