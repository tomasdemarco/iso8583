// Package padding provides functionalities for handling padding in ISO 8583 fields.
package padding

import (
	"fmt"
	"github.com/tomasdemarco/iso8583/encoding"
	"strings"
)

// FillPadder implements the Padder interface for fill padding.
// It pads a field value with a specified character to reach a target length.
type FillPadder struct {
	left bool
	char string
}

// FILL provides pre-configured FillPadder instances for common padding positions.
var FILL = Padders{
	LEFT:  &FillPadder{left: true},
	RIGHT: &FillPadder{left: false},
}

// NewFillPadder creates a new FillPadder.
// `left` indicates if padding should be applied to the left (true) or right (false).
// `char` is the character used for padding.
func NewFillPadder(left bool, char string) Padder {
	return &FillPadder{left, char}
}

// Pack calculates the left and right padding strings for fill padding.
// It returns the padding strings and an error if the value is too long for the field.
func (p *FillPadder) Pack(lengthPackager int, lengthValue int, encoder encoding.Encoder) (string, string, error) {
	if lengthPackager < lengthValue {
		return "", "", fmt.Errorf("%w: value %d, max %d", ErrValueTooLong, lengthValue, lengthPackager)
	}
	if p.left {
		return strings.Repeat(p.char, lengthPackager-lengthValue), "", nil
	}
	return "", strings.Repeat(p.char, lengthPackager-lengthValue), nil
}

// Unpack for FillPadder always returns 0, 0 as fill padding is removed by simply slicing.
func (p *FillPadder) Unpack(_ int) (int, int) {
	return 0, 0
}

func (p *FillPadder) Type() Type {
	return Fill
}

func (p *FillPadder) SetChar(char string) {
	p.char = char
}
