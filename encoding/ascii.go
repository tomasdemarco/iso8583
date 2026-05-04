// Package encoding provides various data encoding and decoding functionalities for ISO 8583 fields.
package encoding

import "fmt"

// AsciiEncoder implements the Encoder interface for ASCII encoding.
// It encodes and decodes strings directly to/from byte slices.
type AsciiEncoder struct {
	length int
}

var ASCII = AsciiEncoder{}

// NewAsciiEncoder creates a new ASCII encoder.
func NewAsciiEncoder() Encoder {
	return &AsciiEncoder{}
}

// Encode converts a string to an ASCII byte slice.
func (e *AsciiEncoder) Encode(src string) ([]byte, error) {
	e.length = len([]byte(src))
	return []byte(src), nil
}

// Decode converts an ASCII byte slice to a string.
// It reads up to the configured length.
func (e *AsciiEncoder) Decode(src []byte) (string, error) {
	if len(src) < e.length {
		return "", fmt.Errorf("%w: expected %d, got %d", ErrNotEnoughDataToDecode, e.length, len(src))
	}
	return string(src[:e.length]), nil
}

// SetLength sets the length for the ASCII encoder.
func (e *AsciiEncoder) SetLength(length int) {
	e.length = length
}

func (e *AsciiEncoder) GetLength() int {
	return e.length
}

func (e *AsciiEncoder) GetType() Encoding {
	return Ascii
}
