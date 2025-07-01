package padding

import "github.com/tomasdemarco/iso8583/encoding"

type Padder interface {
	EncodePad(lengthPackager int, lengthValue int, encoder encoding.Encoder) (string, string, error)
	DecodePad(lengthField int) (int, int)
	SetChar(char string)
	GetChar() string
}

type Padders struct {
	NONE  Padder
	LEFT  Padder
	RIGHT Padder
}
