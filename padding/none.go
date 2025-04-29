package padding

import (
	"github.com/tomasdemarco/iso8583/encoding"
	"github.com/tomasdemarco/iso8583/utils"
)

type NonePadder struct {
	left    bool
	char    utils.ByteFromString
	encoder encoding.Encoder
}

var NONE = Padders{
	NONE: &NonePadder{},
}

func (p *NonePadder) EncodePad(lengthPackager int, lengthValue int) (string, string) {
	return "", ""
}

func (p *NonePadder) DecodePad(lengthField int) (int, int) {
	return 0, 0
}

func (p *NonePadder) SetChar(char utils.ByteFromString) {
	p.char = char
}

func (p *NonePadder) SetEncoder(encoder encoding.Encoder) {
	p.encoder = encoder
}
