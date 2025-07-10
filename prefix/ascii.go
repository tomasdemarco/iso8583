package prefix

import (
	"fmt"
	"github.com/tomasdemarco/iso8583/encoding"
)

// AsciiPrefixer implements the Prefixer interface for ASCII length encoding.
// It encodes and decodes lengths as ASCII digits.
type AsciiPrefixer struct {
	nDigits     int
	encoder     encoding.Encoder
	hex         bool
	isInclusive bool
}

// ASCII provides pre-configured AsciiPrefixer instances for common length types.
var ASCII = Prefixers{
	L:      &AsciiPrefixer{1, &encoding.ASCII{}, false, false},
	LL:     &AsciiPrefixer{2, &encoding.ASCII{}, false, false},
	LLL:    &AsciiPrefixer{3, &encoding.ASCII{}, false, false},
	LLLL:   &AsciiPrefixer{4, &encoding.ASCII{}, false, false},
	LLLLL:  &AsciiPrefixer{5, &encoding.ASCII{}, false, false},
	LLLLLL: &AsciiPrefixer{6, &encoding.ASCII{}, false, false},
}

// NewAsciiPrefixer creates a new AsciiPrefixer with the specified number of digits.
// The `hex` parameter indicates if the length should be treated as hexadecimal.
// The `isInclusive` parameter indicates if the encoded length includes the prefix's own length.
func NewAsciiPrefixer(nDigits int, hex, isInclusive bool) Prefixer {
	return &AsciiPrefixer{nDigits, &encoding.ASCII{}, hex, isInclusive}
}

// EncodeLength encodes the given integer length into an ASCII byte slice.
// It returns the encoded length as a byte slice and an error if encoding fails
// or if the length exceeds the maximum allowed for the configured number of digits.
func (p *AsciiPrefixer) EncodeLength(length int) ([]byte, error) {
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

// DecodeLength decodes an ASCII length from the provided byte slice starting at the given offset.
// It returns the decoded integer length and an error if decoding fails.
func (p *AsciiPrefixer) DecodeLength(b []byte, offset int) (int, error) {
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

// GetPackedLength returns the number of ASCII digits used to encode the length.
func (p *AsciiPrefixer) GetPackedLength() int {
	return p.nDigits
}
