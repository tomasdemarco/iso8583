package encoding

func AsciiDecode(src []byte) (string, error) {
	return string(src), nil
}

func AsciiEncode(src string) []byte {
	return []byte(src)
}
