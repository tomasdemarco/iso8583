package prefix

import (
	"fmt"
	"github.com/tomasdemarco/iso8583/encoding"
)

// EbcdicPrefixer implements the Prefixer interface for EBCDIC length encoding.
type EbcdicPrefixer struct {
	nDigits int
	encoder encoding.Encoder
	hex     bool
}

var EBCDIC = Prefixers{
	L:      &EbcdicPrefixer{1, &encoding.EBCDIC{}, false},
	LL:     &EbcdicPrefixer{2, &encoding.EBCDIC{}, false},
	LLL:    &EbcdicPrefixer{3, &encoding.EBCDIC{}, false},
	LLLL:   &EbcdicPrefixer{4, &encoding.EBCDIC{}, false},
	LLLLL:  &EbcdicPrefixer{5, &encoding.EBCDIC{}, false},
	LLLLLL: &EbcdicPrefixer{6, &encoding.EBCDIC{}, false},
}

// NewEbcdicPrefixer creates a new EbcdicPrefixer with the specified number of digits.
func NewEbcdicPrefixer(nDigits int) EbcdicPrefixer {
	return EbcdicPrefixer{nDigits: nDigits}
}

// EncodeLength encodes the length into the byte slice using EBCDIC.
func (p *EbcdicPrefixer) EncodeLength(length int) ([]byte, error) {
	length, err := lengthInt(length, p.hex)
	if err != nil {
		return nil, err
	}

	return p.encoder.Encode(fmt.Sprintf("%0*d", p.nDigits, length))
}

// DecodeLength decodes the length from the byte slice using EBCDIC.
func (p *EbcdicPrefixer) DecodeLength(b []byte, offset int) (int, error) {
	p.encoder.SetLength(p.nDigits)

	lengthString, err := p.encoder.Decode(b[offset:])
	if err != nil {
		return 0, err
	}

	return lengthStringToInt(lengthString, p.hex)
}

// GetPackedLength returns the number of digits used to encode the length.
func (p *EbcdicPrefixer) GetPackedLength() int {
	return p.nDigits
}

func (p *EbcdicPrefixer) SetHex() {
	p.hex = true
}
