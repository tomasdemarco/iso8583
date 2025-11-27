package header

import (
	"fmt"
	"io"
)

type Header interface {
	Get() any
	Set(any)
	Log() string
}

type BytesHeader struct {
	val []byte
}

func (h *BytesHeader) Get() any       { return h.val }
func (h *BytesHeader) Set(header any) { h.val = header.([]byte) }
func (h *BytesHeader) Log() string    { return fmt.Sprintf("%X", h.val) }

type PackFunc func(val Header) (valueRaw []byte, length int, err error)
type UnpackFunc func(r io.Reader) (val Header, length int, err error)

func Pack(Header) ([]byte, int, error) {
	//not implemented
	return []byte{}, 0, nil
}

func Unpack(r io.Reader) (Header, int, error) {
	//not implemented
	return nil, 0, nil
}
