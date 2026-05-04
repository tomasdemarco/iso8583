// Package encoding provides various data encoding and decoding functionalities for ISO 8583 fields.
package encoding

import (
	"bytes"
	"fmt"
	"strconv"
)

// BCD implements the Encoder interface for BCD (Binary-Coded Decimal) encoding.
// It encodes and decodes decimal strings to/from BCD byte slices.
type BcdEncoder struct {
	length   int
	padRight bool
	odd      bool
}

var BCD = BcdEncoder{}

// NewBcdEncoder creates a new BCD encoder.
// `padRight` indicates whether to right-pad the input string with a '0' if its length is odd
// before encoding to ensure an even number of digits for BCD conversion.
func NewBcdEncoder(padRight bool) Encoder {
	return &BcdEncoder{padRight: padRight}
}

func (e *BcdEncoder) Decode(src []byte) (string, error) {
	var dst bytes.Buffer
	for _, b := range src[:e.length] {
		high := b >> 4
		low := b & 0x0F

		if high > 9 || low > 9 {
			return "", fmt.Errorf("invalid BCD byte: %x", b)
		}

		dst.WriteString(strconv.Itoa(int(high)))
		dst.WriteString(strconv.Itoa(int(low)))
	}

	str := dst.String()
	if e.odd {
		if e.padRight {
			str = str[:len(str)-1]
		} else {
			str = str[1:]
		}
	}

	return str, nil
}

// Encode converts a decimal string to a BCD byte slice.
// If `padLeft` is true and the source string has an odd length, it will be left-padded with '0'.
func (e *BcdEncoder) Encode(src string) ([]byte, error) {
	if len(src)%2 != 0 {
		src = "0" + src
	}

	var result bytes.Buffer
	for i := 0; i < len(src); i += 2 {
		high, err := strconv.ParseUint(string(src[i]), 10, 4)
		if err != nil {
			return nil, fmt.Errorf("BCD string invalid digit: %s", string(src[i]))
		}

		low, err := strconv.ParseUint(string(src[i+1]), 10, 4)
		if err != nil {
			return nil, fmt.Errorf("BCD string invalid digit: %s", string(src[i+1]))
		}

		result.WriteByte(byte(high<<4 | low))
	}

	e.length = len(result.Bytes())

	return result.Bytes(), nil
}

// SetLength sets the length for the BCD encoder.
func (e *BcdEncoder) SetLength(length int) {
	e.odd = length%2 != 0
	e.length = (length + 1) / 2
}

func (e *BcdEncoder) GetLength() int {
	return e.length
}

func (e *BcdEncoder) GetType() Encoding {
	return Bcd
}
