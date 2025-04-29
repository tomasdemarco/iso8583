package padding

import (
	"github.com/tomasdemarco/iso8583/encoding"
	"github.com/tomasdemarco/iso8583/utils"
)

type Padder interface {
	EncodePad(lengthPackager int, lengthValue int) (string, string)
	DecodePad(lengthField int) (int, int)
	SetChar(char utils.ByteFromString)
	SetEncoder(encoder encoding.Encoder)
}

type Padders struct {
	NONE  Padder
	LEFT  Padder
	RIGHT Padder
}
