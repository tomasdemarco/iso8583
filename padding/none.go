package padding

import (
	"github.com/tomasdemarco/iso8583/encoding"
)

type NonePadder struct {
	encoder encoding.Encoder
}

var NONE = Padders{
	NONE: &NonePadder{},
}

func (p *NonePadder) EncodePad(char string, lengthPackager int, lengthValue int, encoder encoding.Encoder) (string, string, error) {
	return "", "", nil
}

func (p *NonePadder) DecodePad(lengthField int) (int, int) {
	return 0, 0
}
