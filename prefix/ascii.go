package prefix

import (
	"fmt"
	"github.com/tomasdemarco/iso8583/encoding"
)

// AsciiPrefixer implements the Prefixer interface for ASCII length encoding.
type AsciiPrefixer struct {
	nDigits     int
	encoder     encoding.Encoder
	hex         bool
	isInclusive bool
}

var ASCII = Prefixers{
	L:      &AsciiPrefixer{1, &encoding.ASCII{}, false, false},
	LL:     &AsciiPrefixer{2, &encoding.ASCII{}, false, false},
	LLL:    &AsciiPrefixer{3, &encoding.ASCII{}, false, false},
	LLLL:   &AsciiPrefixer{4, &encoding.ASCII{}, false, false},
	LLLLL:  &AsciiPrefixer{5, &encoding.ASCII{}, false, false},
	LLLLLL: &AsciiPrefixer{6, &encoding.ASCII{}, false, false},
}

// NewAsciiPrefixer creates a new Prefixer with the specified number of digits.
func NewAsciiPrefixer(nDigits int, hex, isInclusive bool) Prefixer {
	return &AsciiPrefixer{nDigits, &encoding.ASCII{}, hex, isInclusive}
}

// EncodeLength encodes the length into the byte slice.
func (p *AsciiPrefixer) EncodeLength(length int) ([]byte, error) {
	length, err := lengthInt(length, p.hex)
	if err != nil {
		return nil, err
	}

	if p.isInclusive {
		length += p.nDigits
	}

	return p.encoder.Encode(fmt.Sprintf("%0*d", p.nDigits, length))
}

// DecodeLength decodes the length from the byte slice.
func (p *AsciiPrefixer) DecodeLength(b []byte, offset int) (int, error) {
	p.encoder.SetLength(p.nDigits)

	lengthString, err := p.encoder.Decode(b[offset:])
	if err != nil {
		return 0, err
	}

	length, err := lengthStringToInt(lengthString, p.hex)
	if err != nil {
		return 0, err
	}

	if p.isInclusive && length >= p.nDigits {
		return length - p.nDigits, nil
	}

	return length, nil
}

// GetPackedLength returns the number of digits used to encode the length.
func (p *AsciiPrefixer) GetPackedLength() int {
	return p.nDigits
}

func (p *AsciiPrefixer) SetHex(hex bool) {
	p.hex = hex
}

func (p *AsciiPrefixer) SetIsInclusive(isInclusive bool) {
	p.isInclusive = isInclusive
}
