package prefix

import (
	"fmt"
	"github.com/tomasdemarco/iso8583/encoding"
	"github.com/tomasdemarco/iso8583/utils"
	"strconv"
)

// BinaryPrefixer implements the Prefixer interface for BINARY length encoding.
type BinaryPrefixer struct {
	nBytes  int
	encoder encoding.Encoder
	hex     bool
}

var BINARY = BinaryPrefixers{
	B:  &BinaryPrefixer{1, &encoding.HEX{}, false},
	BB: &BinaryPrefixer{2, &encoding.HEX{}, false},
}

// NewBinaryPrefixer creates a new BinaryPrefixer with the specified number of bytes.
func NewBinaryPrefixer(nBytes int) BinaryPrefixer {
	return BinaryPrefixer{nBytes: nBytes}
}

// EncodeLength encodes the length into the byte slice.
func (p *BinaryPrefixer) EncodeLength(length int) ([]byte, error) {
	if p.hex {
		length64, err := strconv.ParseInt(fmt.Sprintf("%d", length), 10, 16)
		if err != nil {
			return nil, err
		}

		length = int(length64)
	}

	b, err := p.encoder.Encode(fmt.Sprintf("%d", length))
	if err != nil {
		return nil, err
	}

	return utils.ZeroPadLeft(b, p.nBytes)
}

// DecodeLength decodes the length from the byte slice.
func (p *BinaryPrefixer) DecodeLength(b []byte, offset int) (int, error) {
	p.encoder.SetLength(p.nBytes)

	lengthString, err := p.encoder.Decode(b[offset:])
	if err != nil {
		return 0, err
	}

	length, err := strconv.Atoi(lengthString)
	if err != nil {
		return 0, err
	}

	return length, nil
}

// GetPackedLength returns the number of bytes used to encode the length.
func (p *BinaryPrefixer) GetPackedLength() int {
	return p.nBytes
}

func (p *BinaryPrefixer) SetHex() {
	p.hex = true
}
