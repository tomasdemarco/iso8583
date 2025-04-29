package prefix

import (
	"github.com/tomasdemarco/iso8583/encoding"
)

// NonePrefixer implements the Prefixer interface for NONE length encoding.
type NonePrefixer struct {
	nDigits int
	encoder encoding.Encoder
	hex     bool
}

var NONE = Prefixers{
	Fixed: &NonePrefixer{},
}

// NewNonePrefixer creates a new Prefixer with the specified number of digits.
func NewNonePrefixer(nDigits int) NonePrefixer {
	return NonePrefixer{nDigits: nDigits}
}

// EncodeLength encodes the length into the byte slice.
func (p *NonePrefixer) EncodeLength(length int) ([]byte, error) {
	return nil, nil
}

// DecodeLength decodes the length from the byte slice.
func (p *NonePrefixer) DecodeLength(b []byte, offset int) (int, error) {
	return 0, nil
}

// GetPackedLength returns the number of digits used to encode the length.
func (p *NonePrefixer) GetPackedLength() int {
	return p.nDigits
}

func (p *NonePrefixer) SetHex() {
	p.hex = true
}
