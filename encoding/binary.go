// Package encoding provides various data encoding and decoding functionalities for ISO 8583 fields.
package encoding

import (
	"fmt"
	"github.com/tomasdemarco/iso8583/utils"
)

// BINARY implements the Encoder interface for raw binary encoding.
// It treats input strings as hexadecimal representations of binary data.
type BINARY struct {
	length int
}

// NewBinaryEncoder creates a new BINARY encoder.
func NewBinaryEncoder() Encoder {
	return &BINARY{}
}

// Encode converts a hexadecimal string into a raw binary byte slice.
// If the source string has an odd length, it will be left-padded with '0'.
func (e *BINARY) Encode(src string) ([]byte, error) {
	if len(src)%2 != 0 {
		src = "0" + src
	}
	return utils.Hex2Byte(src), nil
}

// Decode converts a raw binary byte slice into an uppercase hexadecimal string.
// It reads up to the configured length.
func (e *BINARY) Decode(src []byte) (string, error) {
	if len(src) < e.length {
		return "", fmt.Errorf("BINARY decode: not enough data to read. expected %d, got %d", e.length, len(src))
	}
	return fmt.Sprintf("%X", src[:e.length]), nil
}

// SetLength sets the length for the BINARY encoder.
func (e *BINARY) SetLength(length int) {
	e.length = length
}
