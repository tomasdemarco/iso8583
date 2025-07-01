package padding

import (
	"github.com/tomasdemarco/iso8583/encoding"
)

type ParityPadder struct {
	left bool
	char string
}

var PARITY = Padders{
	LEFT:  &ParityPadder{left: true},
	RIGHT: &ParityPadder{left: false},
}

func NewParityPadder(left bool, char string) Padder {
	return &ParityPadder{left, char}
}

func (p *ParityPadder) EncodePad(_, lengthValue int, _ encoding.Encoder) (string, string, error) {
	if lengthValue%2 != 0 {
		if p.left {
			return p.char, "", nil

		}
		return "", p.char, nil
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

func (p *ParityPadder) SetChar(char string) {
	p.char = char
}

func (p *ParityPadder) GetChar() string {
	return p.char
}
