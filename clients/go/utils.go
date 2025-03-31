package scas

func join2Bytes(b1, b2, b3 []byte) []byte {
	return append(b1, append(b2, b3...)...)
}

func joinBytes(b1, b2 []byte) []byte {
	return append(b1, b2...)
}

func joinByte(b1 []byte, b2 byte) []byte {
	return append(b1, b2)
}
