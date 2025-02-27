package padding

import "github.com/tomasdemarco/iso8583/utils"

type Padding struct {
	Type     Type                 `json:"type"`
	Position Position             `json:"position"`
	Char     utils.ByteFromString `json:"char"`
}

func Unpack(padding Padding, lengthField int) (int, int) {
	switch padding.Position {
	case Right:
		pad := RightDecode(padding.Type, lengthField)
		return pad, 0
	case Left:
		pad := LeftDecode(padding.Type, lengthField)
		return 0, pad
	default:
		return 0, 0
	}
}

func Pack(padding Padding, lengthPackager int, lengthValue int) (string, string) {
	switch padding.Position {
	case Right:
		padResult := RightEncode(padding.Type, padding.Char, lengthPackager, lengthValue)
		return padResult, ""
	case Left:
		padResult := LeftEncode(padding.Type, padding.Char, lengthPackager, lengthValue)
		return "", padResult
	default:
		return "", ""
	}
}
