package prefix

import (
	"fmt"
	"github.com/tomasdemarco/iso8583/encoding"
	"strconv"
)

// BcdPrefixer implements the Prefixer interface for BCD length encoding.
type BcdPrefixer struct {
	nDigits int
	encoder encoding.Encoder
	hex     bool
}

var BCD = Prefixers{
	L:      &BcdPrefixer{2, &encoding.BCD{}, false},
	LL:     &BcdPrefixer{2, &encoding.BCD{}, false},
	LLL:    &BcdPrefixer{4, &encoding.BCD{}, false},
	LLLL:   &BcdPrefixer{4, &encoding.BCD{}, false},
	LLLLL:  &BcdPrefixer{6, &encoding.BCD{}, false},
	LLLLLL: &BcdPrefixer{6, &encoding.BCD{}, false},
}

// NewBcdPrefixer creates a new BcdPrefixer with the specified number of digits.
func NewBcdPrefixer(nDigits int) *BcdPrefixer {
	return &BcdPrefixer{nDigits: nDigits}
}

// EncodeLength encodes the length into the byte slice.
func (p *BcdPrefixer) EncodeLength(length int) ([]byte, error) {
	if p.hex {
		length64, err := strconv.ParseInt(fmt.Sprintf("%x", length), 10, 64)
		if err != nil {
			return nil, err
		}

		length = int(length64)
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

	if p.hex {
		length, err := strconv.ParseInt(lengthString, 16, 10)
		if err != nil {
			return 0, err
		}
		return int(length), nil
	} else {
		length, err := strconv.Atoi(lengthString)
		if err != nil {
			return 0, err
		}
		return length, nil
	}
}

// GetPackedLength returns the number of digits used to encode the length.
func (p *BcdPrefixer) GetPackedLength() int {
	return p.nDigits / 2
}

func (p *BcdPrefixer) SetHex() {
	p.hex = true
}
