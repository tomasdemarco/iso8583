// Package prefix provides functionalities for handling length prefixes in ISO 8583 messages.
package prefix

import (
	"fmt"
	"github.com/tomasdemarco/iso8583/encoding"
)

// BinaryPrefixer implements the Prefixer interface for BINARY length encoding.
// It encodes and decodes lengths as binary bytes.
type BinaryPrefixer struct {
	nBytes      int
	encoder     encoding.Encoder
	hex         bool
	isInclusive bool
}

// BINARY provides pre-configured BinaryPrefixer instances for common binary length types.
var BINARY = BinaryPrefixers{
	B:  &BinaryPrefixer{1, &encoding.BINARY{}, true, false},
	BB: &BinaryPrefixer{2, &encoding.BINARY{}, true, false},
}

// NewBinaryPrefixer creates a new BinaryPrefixer with the specified number of bytes.
// The `isInclusive` parameter indicates if the encoded length includes the prefix's own length.
func NewBinaryPrefixer(nBytes int, isInclusive bool) Prefixer {
	return &BinaryPrefixer{(nBytes + 1) / 2, &encoding.BINARY{}, true, isInclusive}
}

// EncodeLength encodes the given integer length into a binary byte slice.
// It returns the encoded length as a byte slice and an error if encoding fails
// or if the length exceeds the maximum allowed for the configured number of bytes.
func (p *BinaryPrefixer) EncodeLength(length int) ([]byte, error) {
	err := validateMaxLimit(length, p.nBytes*2, p.hex)
	if err != nil {
		return nil, err
	}

	if p.isInclusive {
		length += p.nBytes
	}

	lenStr := intToLenStr(length, p.hex)

	return p.encoder.Encode(fmt.Sprintf("%0*s", p.nBytes*2, lenStr))
}

// DecodeLength decodes a binary length from the provided byte slice starting at the given offset.
// It returns the decoded integer length and an error if decoding fails.
func (p *BinaryPrefixer) DecodeLength(b []byte, offset int) (int, error) {
	p.encoder.SetLength(p.nBytes)

	lengthString, err := p.encoder.Decode(b[offset:])
	if err != nil {
		return 0, fmt.Errorf("%w: %w", ErrFailedToDecodeLength, err)
	}

	length, err := lenStrToInt(lengthString, p.hex)
	if err != nil {
		return 0, fmt.Errorf("%w: %w", ErrInvalidLengthStringConversion, err)
	}

	if p.isInclusive {
		if length < p.nBytes {
			return 0, fmt.Errorf("%w: decoded length %d is less than prefix length %d", ErrInvalidLengthStringConversion, length, p.nBytes)
		}
		return length - p.nBytes, nil
	}

	return length, nil
}

// GetPackedLength returns the number of bytes used to encode the length.
func (p *BinaryPrefixer) GetPackedLength() int {
	return p.nBytes
}
