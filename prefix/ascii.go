package prefix

import (
	"fmt"
	"github.com/tomasdemarco/iso8583/encoding"
	"strconv"
)

// AsciiPrefixer implements the Prefixer interface for ASCII length encoding.
type AsciiPrefixer struct {
	nDigits int
	encoder encoding.Encoder
	hex     bool
}

var ASCII = Prefixers{
	L:      &AsciiPrefixer{1, &encoding.ASCII{}, false},
	LL:     &AsciiPrefixer{2, &encoding.ASCII{}, false},
	LLL:    &AsciiPrefixer{3, &encoding.ASCII{}, false},
	LLLL:   &AsciiPrefixer{4, &encoding.ASCII{}, false},
	LLLLL:  &AsciiPrefixer{5, &encoding.ASCII{}, false},
	LLLLLL: &AsciiPrefixer{6, &encoding.ASCII{}, false},
}

// NewAsciiPrefixer creates a new Prefixer with the specified number of digits.
func NewAsciiPrefixer(nDigits int) AsciiPrefixer {
	return AsciiPrefixer{nDigits: nDigits}
}

// EncodeLength encodes the length into the byte slice.
func (p *AsciiPrefixer) EncodeLength(length int) ([]byte, error) {
	if p.hex {
		length64, err := strconv.ParseInt(fmt.Sprintf("%d", length), 10, 16)
		if err != nil {
			return nil, err
		}

		length = int(length64)
	}
	return p.encoder.Encode(fmt.Sprintf("%0*d", p.nDigits, length))
}

// DecodeLength decodes the length from the byte slice.
func (p *AsciiPrefixer) DecodeLength(b []byte, offset int) (int, error) {
	p.encoder.SetLength(p.nDigits)

	lengthString, err := p.encoder.Decode(b[offset:offset])
	if err != nil {
		return 0, err
	}

	length, err := strconv.Atoi(lengthString)
	if err != nil {
		return 0, err
	}

	return length, nil
}

// GetPackedLength returns the number of digits used to encode the length.
func (p *AsciiPrefixer) GetPackedLength() int {
	return p.nDigits
}

func (p *AsciiPrefixer) SetHex() {
	p.hex = true
}
