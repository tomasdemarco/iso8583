// Package prefix provides functionalities for handling length prefixes in ISO 8583 messages.
package prefix

import (
	"fmt"
	"github.com/tomasdemarco/iso8583/encoding"
)

// BcdPrefixer implements the Prefixer interface for BCD length encoding.
// It encodes and decodes lengths as BCD (Binary-Coded Decimal) digits.
type BcdPrefixer struct {
	nDigits     int
	encoder     encoding.Encoder
	hex         bool
	isInclusive bool
}

// BCD provides pre-configured BcdPrefixer instances for common length types.
var BCD = Prefixers{
	L:      &BcdPrefixer{2, &encoding.BCD{}, false, false},
	LL:     &BcdPrefixer{2, &encoding.BCD{}, false, false},
	LLL:    &BcdPrefixer{4, &encoding.BCD{}, false, false},
	LLLL:   &BcdPrefixer{4, &encoding.BCD{}, false, false},
	LLLLL:  &BcdPrefixer{6, &encoding.BCD{}, false, false},
	LLLLLL: &BcdPrefixer{6, &encoding.BCD{}, false, false},
}

// NewBcdPrefixer creates a new BcdPrefixer with the specified number of digits.
// The `hex` parameter indicates if the length should be treated as hexadecimal.
// The `isInclusive` parameter indicates if the encoded length includes the prefix's own length.
func NewBcdPrefixer(nDigits int, hex, isInclusive bool) Prefixer {
	return &BcdPrefixer{nDigits, encoding.NewBcdEncoder(true), hex, isInclusive}
}

// EncodeLength encodes the given integer length into a BCD byte slice.
// It returns the encoded length as a byte slice and an error if encoding fails
// or if the length exceeds the maximum allowed for the configured number of digits.
func (p *BcdPrefixer) EncodeLength(length int) ([]byte, error) {
	err := validateMaxLimit(length, p.nDigits, p.hex)
	if err != nil {
		return nil, err
	}

	if p.isInclusive {
		length += p.nDigits
	}

	lenStr := intToLenStr(length, p.hex)

	return p.encoder.Encode(fmt.Sprintf("%0*s", p.nDigits, lenStr))
}

// DecodeLength decodes a BCD length from the provided byte slice starting at the given offset.
// It returns the decoded integer length and an error if decoding fails.
func (p *BcdPrefixer) DecodeLength(b []byte, offset int) (int, error) {
	p.encoder.SetLength(p.GetPackedLength())

	lengthString, err := p.encoder.Decode(b[offset:])
	if err != nil {
		return 0, err
	}

	length, err := lenStrToInt(lengthString, p.hex)
	if err != nil {
		return 0, err
	}

	if p.isInclusive && length >= p.nDigits {
		return length - p.nDigits, nil
	}

	return length, nil
}

// GetPackedLength returns the number of bytes used to encode the length.
// For BCD, this is (nDigits + 1) / 2.
func (p *BcdPrefixer) GetPackedLength() int {
	return (p.nDigits + 1) / 2
}
