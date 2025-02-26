package padding

type Padding struct {
	Type     Type     `json:"type"`
	Position Position `json:"position"`
	Char     byte     `json:"char"`
}

func Unpack(padding Padding) (int, int) {
	switch padding.Position {
	case Right:
		pad := RightDecode(padding.Type)
		return pad, 0
	case Left:
		pad := LeftDecode(padding.Type)
		return 0, pad
	default:
		return 0, 0
	}
}

func Pack(padding Padding, lengthField int, lengthValue int) (string, string) {
	switch padding.Position {
	case Right:
		padResult := RightEncode(padding.Type, padding.Char, lengthValue, lengthField)
		return padResult, ""
	case Left:
		padResult := LeftEncode(padding.Type, padding.Char, lengthValue, lengthField)
		return "", padResult
	default:
		return "", ""
	}
}
