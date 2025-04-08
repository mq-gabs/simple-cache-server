package main

func joinBytes(bytes [][]byte) []byte {
	var res []byte

	for _, b := range bytes {
		res = append(res, b...)
	}

	return res
}
