package padding

import (
	"github.com/tomasdemarco/iso8583/encoding"
)

type ParityPadder struct {
	left bool
}

var PARITY = Padders{
	LEFT:  &ParityPadder{left: true},
	RIGHT: &ParityPadder{left: false},
}

func (p *ParityPadder) EncodePad(char string, lengthPackager int, lengthValue int, encoder encoding.Encoder) (string, string, error) {
	if lengthValue%2 != 0 {
		if p.left {
			return string(char), "", nil

		}
		return "", string(char), nil
	}
	return "", "", nil
}

func (p *ParityPadder) DecodePad(lengthField int) (int, int) {
	if lengthField%2 == 0 {
		return 0, 0
	}

	if p.left {
		return 1, 0
	}
	return 0, 1
}
