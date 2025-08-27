package utils

import (
	"bytes"
	"testing"
)

func TestJoinBytes(t *testing.T) {
	value1 := []byte{0x00}
	value2 := []byte{0x01}
	var sep byte = 0x33
	expected := []byte{value1[0], sep, value2[0]}

	res := JoinBytes(sep, value1, value2)

	if !bytes.Equal(res, expected) {
		t.Fatalf("Bytes aren't equal: expected: %v, got: %v", expected, res)
	}
}
