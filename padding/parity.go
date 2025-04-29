package padding

import (
	"github.com/tomasdemarco/iso8583/encoding"
	"github.com/tomasdemarco/iso8583/utils"
)

type ParityPadder struct {
	left    bool
	char    utils.ByteFromString
	encoder encoding.Encoder
}

var PARITY = Padders{
	LEFT:  &ParityPadder{left: true},
	RIGHT: &ParityPadder{left: false},
}

func (p *ParityPadder) EncodePad(lengthPackager int, lengthValue int) (string, string) {
	if lengthValue%2 != 0 {
		if p.left {
			return string(p.char), ""

		}
		return "", string(p.char)
	}
	return "", ""
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

func (p *ParityPadder) SetChar(char utils.ByteFromString) {
	p.char = char
}

func (p *ParityPadder) SetEncoder(encoder encoding.Encoder) {
	p.encoder = encoder
}
