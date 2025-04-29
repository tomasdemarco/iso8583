package padding

import (
	"fmt"
	"github.com/tomasdemarco/iso8583/encoding"
	"github.com/tomasdemarco/iso8583/utils"
	"strings"
)

type FillPadder struct {
	left    bool
	char    utils.ByteFromString
	encoder encoding.Encoder
}

var FILL = Padders{
	LEFT:  &FillPadder{left: true},
	RIGHT: &FillPadder{left: false},
}

func (p *FillPadder) EncodePad(lengthPackager int, lengthValue int) (string, string) {
	fmt.Println(lengthPackager, lengthValue)
	if _, ok := p.encoder.(*encoding.BCD); ok {
		lengthPackager = lengthPackager * 2
	}
	fmt.Println(lengthPackager, lengthValue)
	if p.left {
		return strings.Repeat(string(p.char), lengthPackager-lengthValue), ""
	}
	return "", strings.Repeat(string(p.char), lengthPackager-lengthValue)
}

func (p *FillPadder) DecodePad(lengthField int) (int, int) {
	return 0, 0
}

func (p *FillPadder) SetChar(char utils.ByteFromString) {
	p.char = char
}

func (p *FillPadder) SetEncoder(encoder encoding.Encoder) {
	p.encoder = encoder
}
