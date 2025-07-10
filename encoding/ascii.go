// Package encoding provides various data encoding and decoding functionalities for ISO 8583 fields.
package encoding

import "fmt"

// ASCII implements the Encoder interface for ASCII encoding.
// It encodes and decodes strings directly to/from byte slices.
type ASCII struct {
	length int
}

// NewAsciiEncoder creates a new ASCII encoder.
func NewAsciiEncoder() Encoder {
	return &ASCII{}
}

// Encode converts a string to an ASCII byte slice.
func (e *ASCII) Encode(src string) ([]byte, error) {
	return []byte(src), nil
}

// Decode converts an ASCII byte slice to a string.
// It reads up to the configured length.
func (e *ASCII) Decode(src []byte) (string, error) {
	if len(src) < e.length {
		return "", fmt.Errorf("ASCII decode: not enough data to read. expected %d, got %d", e.length, len(src))
	}
	return string(src[:e.length]), nil
}

// SetLength sets the length for the ASCII encoder.
func (e *ASCII) SetLength(length int) {
	e.length = length
}
