package padding

import (
	"github.com/tomasdemarco/iso8583/encoding"
)

type NonePadder struct {
	encoder encoding.Encoder
	char    string
}

var NONE = Padders{
	NONE: &NonePadder{},
}

func (p *NonePadder) EncodePad(_ int, _ int, _ encoding.Encoder) (string, string, error) {
	return "", "", nil
}

func (p *NonePadder) DecodePad(_ int) (int, int) {
	return 0, 0
}

func (p *NonePadder) SetChar(char string) {
	p.char = char
}

func (p *NonePadder) GetChar() string {
	return p.char
}
