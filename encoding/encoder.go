package encoding

type Encoder interface {
	Encode(src string) ([]byte, error)
	Decode(src []byte) (string, error)
	SetLength(length int)
}
