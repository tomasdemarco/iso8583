package padding

import (
	"fmt"
	"github.com/tomasdemarco/iso8583/encoding"
	"strings"
)

type FillPadder struct {
	left bool
}

var FILL = Padders{
	LEFT:  &FillPadder{left: true},
	RIGHT: &FillPadder{left: false},
}

func (p *FillPadder) EncodePad(char string, lengthPackager int, lengthValue int, encoder encoding.Encoder) (string, string, error) {
	if _, ok := encoder.(*encoding.BCD); ok {
		lengthPackager = lengthPackager * 2
	}
	if lengthPackager < lengthValue {
		return "", "", fmt.Errorf("value %d too long, maximum %d", lengthValue, lengthPackager)
	}
	if p.left {
		return strings.Repeat(char, lengthPackager-lengthValue), "", nil
	}
	return "", strings.Repeat(char, lengthPackager-lengthValue), nil
}

func (p *FillPadder) DecodePad(lengthField int) (int, int) {
	return 0, 0
}
