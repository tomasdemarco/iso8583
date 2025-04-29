package encoding

// Encoder implements the Encoder interface for ASCII encoding.
type ASCII struct {
	length int
}

func NewAsciiEncoder() ASCII {
	return ASCII{}
}

func (e *ASCII) Encode(src string) ([]byte, error) {
	return []byte(src), nil
}

func (e *ASCII) Decode(src []byte) (string, error) {
	return string(src[:e.length]), nil
}

func (e *ASCII) SetLength(length int) {
	e.length = length
}
