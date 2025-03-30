package scas

func Join2Bytes(b1, b2, b3 []byte) []byte {
	return append(b1, append(b2, b3...)...)
}

func JoinBytes(b1, b2 []byte) []byte {
	return append(b1, b2...)
}

func JoinByte(b1 []byte, b2 byte) []byte {
	return append(b1, b2)
}
