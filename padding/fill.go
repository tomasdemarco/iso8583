package padding

import (
	"fmt"
	"github.com/tomasdemarco/iso8583/encoding"
	"strings"
)

type FillPadder struct {
	left bool
	char string
}

var FILL = Padders{
	LEFT:  &FillPadder{left: true},
	RIGHT: &FillPadder{left: false},
}

func NewFillPadder(left bool, char string) Padder {
	return &FillPadder{left, char}
}

func (p *FillPadder) EncodePad(lengthPackager int, lengthValue int, encoder encoding.Encoder) (string, string, error) {
	if _, ok := encoder.(*encoding.BCD); ok {
		lengthPackager = lengthPackager * 2
	}
	if lengthPackager < lengthValue {
		return "", "", fmt.Errorf("value %d too long, maximum %d", lengthValue, lengthPackager)
	}
	if p.left {
		return strings.Repeat(p.char, lengthPackager-lengthValue), "", nil
	}
	return "", strings.Repeat(p.char, lengthPackager-lengthValue), nil
}

func (p *FillPadder) DecodePad(_ int) (int, int) {
	return 0, 0
}

func (p *FillPadder) SetChar(char string) {
	p.char = char
}

func (p *FillPadder) GetChar() string {
	return p.char
}
