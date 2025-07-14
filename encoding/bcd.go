// Package encoding provides various data encoding and decoding functionalities for ISO 8583 fields.
package encoding

import (
	"bytes"
	"fmt"
)

// BCD implements the Encoder interface for BCD (Binary-Coded Decimal) encoding.
// It encodes and decodes decimal strings to/from BCD byte slices.
type BCD struct {
	length  int
	padLeft bool
}

// NewBcdEncoder creates a new BCD encoder.
// `padLeft` indicates whether to left-pad the input string with a '0' if its length is odd
// before encoding to ensure an even number of digits for BCD conversion.
func NewBcdEncoder(padLeft bool) Encoder {
	return &BCD{padLeft: padLeft}
}

// Encode converts a decimal string to a BCD byte slice.
// If `padLeft` is true and the source string has an odd length, it will be left-padded with '0'.
func (e *BCD) Encode(src string) ([]byte, error) {
	start := 0
	d := make([]byte, (len(src)+1)/2)

	if len(src)%2 == 1 && e.padLeft {
		start = 1
	}

	for i := start; i < len(src)+start; i++ {
		n := i / 2
		digit := src[i-start] - '0'
		if i%2 == 1 {
			d[n] |= digit
		} else {
			d[n] |= digit << 4
		}
	}
	return d, nil
}

// Decode converts a BCD byte slice to a decimal string.
// It reads up to the configured length.
func (e *BCD) Decode(src []byte) (string, error) {
	if len(src) < e.length {
		return "", fmt.Errorf("%w: expected %d, got %d", ErrNotEnoughDataToDecode, e.length, len(src))
	}

	src = src[:e.length]
	start := 0
	var d bytes.Buffer

	for i := start; i < len(src)*2+start; i++ {
		shift := 0
		if i%2 == 1 {
			shift = 0
		} else {
			shift = 4
		}

		c := (src[i/2] >> shift) & 0xF
		var char rune
		if c < 10 {
			char = rune(c + '0')
		} else {
			char = rune(c - 10 + 'A')
		}

		if char == 'D' {
			char = '='
		}
		d.WriteRune(char)
	}
	return d.String(), nil
}

// SetLength sets the length for the BCD encoder.
func (e *BCD) SetLength(length int) {
	e.length = length
}
