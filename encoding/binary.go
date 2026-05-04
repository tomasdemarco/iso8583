// Package encoding provides various data encoding and decoding functionalities for ISO 8583 fields.
package encoding

import (
	"fmt"
	"github.com/tomasdemarco/iso8583/utils"
)

// BinaryEncoder implements the Encoder interface for raw binary encoding.
// It treats input strings as hexadecimal representations of binary data.
type BinaryEncoder struct {
	length int
}

var BINARY = BinaryEncoder{}

// NewBinaryEncoder creates a new BinaryEncoder encoder.
func NewBinaryEncoder() Encoder {
	return &BinaryEncoder{}
}

// Encode converts a hexadecimal string into a raw binary byte slice.
// If the source string has an odd length, it will be left-padded with '0'.
func (e *BinaryEncoder) Encode(src string) ([]byte, error) {
	if len(src)%2 != 0 {
		src = "0" + src
	}
	dst := utils.Hex2Byte(src)
	e.length = len(dst)
	return dst, nil
}

// Decode converts a raw binary byte slice into an uppercase hexadecimal string.
// It reads up to the configured length.
func (e *BinaryEncoder) Decode(src []byte) (string, error) {
	if len(src) < e.length {
		return "", fmt.Errorf("%w: expected %d, got %d", ErrNotEnoughDataToDecode, e.length, len(src))
	}
	return fmt.Sprintf("%X", src[:e.length]), nil
}

// SetLength sets the length for the BinaryEncoder encoder.
func (e *BinaryEncoder) SetLength(length int) {
	e.length = length
}

func (e *BinaryEncoder) GetLength() int {
	return e.length
}

func (e *BinaryEncoder) GetType() Encoding {
	return Binary
}
