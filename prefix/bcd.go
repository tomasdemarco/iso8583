package prefix

import (
	"fmt"
	"github.com/tomasdemarco/iso8583/encoding"
)

// BcdPrefixer implements the Prefixer interface for BCD length encoding.
type BcdPrefixer struct {
	nDigits     int
	encoder     encoding.Encoder
	hex         bool
	isInclusive bool
}

var BCD = Prefixers{
	L:      &BcdPrefixer{2, &encoding.BCD{}, false, false},
	LL:     &BcdPrefixer{2, &encoding.BCD{}, false, false},
	LLL:    &BcdPrefixer{4, &encoding.BCD{}, false, false},
	LLLL:   &BcdPrefixer{4, &encoding.BCD{}, false, false},
	LLLLL:  &BcdPrefixer{6, &encoding.BCD{}, false, false},
	LLLLLL: &BcdPrefixer{6, &encoding.BCD{}, false, false},
}

// NewBcdPrefixer creates a new BcdPrefixer with the specified number of digits.
func NewBcdPrefixer(nDigits int, hex, isInclusive bool) Prefixer {
	return &BcdPrefixer{nDigits, &encoding.BCD{}, hex, isInclusive}
}

// EncodeLength encodes the length into the byte slice.
func (p *BcdPrefixer) EncodeLength(length int) ([]byte, error) {
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
func (p *BcdPrefixer) DecodeLength(b []byte, offset int) (int, error) {
	p.encoder.SetLength(p.GetPackedLength())

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
func (p *BcdPrefixer) GetPackedLength() int {
	return p.nDigits / 2
}

func (p *BcdPrefixer) SetHex(hex bool) {
	p.hex = hex
}

func (p *BcdPrefixer) SetIsInclusive(isInclusive bool) {
	p.isInclusive = isInclusive
}
