// Package prefix provides functionalities for handling length prefixes in ISO 8583 messages.
package prefix

import (
	"fmt"
	"github.com/tomasdemarco/iso8583/encoding"
)

// EbcdicPrefixer implements the Prefixer interface for EBCDIC length encoding.
// It encodes and decodes lengths as EBCDIC digits.
type EbcdicPrefixer struct {
	nDigits     int
	encoder     encoding.Encoder
	hex         bool
	isInclusive bool
}

// EBCDIC provides pre-configured EbcdicPrefixer instances for common length types.
var EBCDIC = Prefixers{
	L:      &EbcdicPrefixer{1, &encoding.EBCDIC{}, false, false},
	LL:     &EbcdicPrefixer{2, &encoding.EBCDIC{}, false, false},
	LLL:    &EbcdicPrefixer{3, &encoding.EBCDIC{}, false, false},
	LLLL:   &EbcdicPrefixer{4, &encoding.EBCDIC{}, false, false},
	LLLLL:  &EbcdicPrefixer{5, &encoding.EBCDIC{}, false, false},
	LLLLLL: &EbcdicPrefixer{6, &encoding.EBCDIC{}, false, false},
}

// NewEbcdicPrefixer creates a new EbcdicPrefixer with the specified number of digits.
// The `hex` parameter indicates if the length should be treated as hexadecimal.
// The `isInclusive` parameter indicates if the encoded length includes the prefix's own length.
func NewEbcdicPrefixer(nDigits int, hex, isInclusive bool) Prefixer {
	return &EbcdicPrefixer{nDigits, &encoding.EBCDIC{}, hex, isInclusive}
}

// EncodeLength encodes the given integer length into an EBCDIC byte slice.
// It returns the encoded length as a byte slice and an error if encoding fails
// or if the length exceeds the maximum allowed for the configured number of digits.
func (p *EbcdicPrefixer) EncodeLength(length int) ([]byte, error) {
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

// DecodeLength decodes an EBCDIC length from the provided byte slice starting at the given offset.
// It returns the decoded integer length and an error if decoding fails.
func (p *EbcdicPrefixer) DecodeLength(b []byte, offset int) (int, error) {
	p.encoder.SetLength(p.nDigits)

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

// GetPackedLength returns the number of EBCDIC digits used to encode the length.
func (p *EbcdicPrefixer) GetPackedLength() int {
	return p.nDigits
}
