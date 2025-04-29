package encoding

import (
	"fmt"
	"github.com/tomasdemarco/iso8583/utils"
)

// BINARY implements the Encoder interface for BINARY encoding.
type BINARY struct {
	length int
}

func NewBinaryEncoder() BINARY {
	return BINARY{}
}

func (e *BINARY) Encode(src string) ([]byte, error) {
	return utils.Hex2Byte(src), nil
}

func (e *BINARY) Decode(src []byte) (string, error) {
	return fmt.Sprintf("%x", src[:e.length]), nil
}

func (e *BINARY) SetLength(length int) {
	e.length = length
}
