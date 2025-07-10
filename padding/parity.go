// Package padding provides functionalities for handling padding in ISO 8583 fields.
package padding

import (
	"github.com/tomasdemarco/iso8583/encoding"
)

// ParityPadder implements the Padder interface for parity padding.
// It adds a single padding character to make the field length even if it's odd.
type ParityPadder struct {
	left bool
	char string
}

// PARITY provides pre-configured ParityPadder instances for common padding positions.
var PARITY = Padders{
	LEFT:  &ParityPadder{left: true},
	RIGHT: &ParityPadder{left: false},
}

// NewParityPadder creates a new ParityPadder.
// `left` indicates if padding should be applied to the left (true) or right (false).
// `char` is the character used for padding.
func NewParityPadder(left bool, char string) Padder {
	return &ParityPadder{left, char}
}

// EncodePad calculates the left and right padding strings for parity padding.
// It adds one padding character if the value length is odd.
func (p *ParityPadder) EncodePad(_, lengthValue int, _ encoding.Encoder) (string, string, error) {
	if lengthValue%2 != 0 {
		if p.left {
			return p.char, "", nil

		}
		return "", p.char, nil
	}
	return "", "", nil
}

// DecodePad calculates the number of characters to remove for parity padding.
// It returns 1 for the appropriate side if the field length is odd, otherwise 0.
func (p *ParityPadder) DecodePad(lengthField int) (int, int) {
	if lengthField%2 == 0 {
		return 0, 0
	}

	if p.left {
		return 1, 0
	}
	return 0, 1
}

// SetChar sets the padding character for the ParityPadder.
func (p *ParityPadder) SetChar(char string) {
	p.char = char
}

// GetChar returns the padding character used by the ParityPadder.
func (p *ParityPadder) GetChar() string {
	return p.char
}
