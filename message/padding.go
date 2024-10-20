package message

import (
	"iso8583/padding"
)

func (m *Message) UnpackPadding(field string) (int, int) {
	switch m.Packager.Fields[field].Padding.Position {
	case "RIGHT":
		pad := padding.RightDecode(m.Packager.Fields[field].Padding.Type)
		return pad, 0
	case "LEFT":
		pad := padding.LeftDecode(m.Packager.Fields[field].Padding.Type)
		return 0, pad
	default:
		return 0, 0
	}
}

func (m *Message) PackPadding(field string) (string, string) {
	typePad := m.Packager.Fields[field].Padding.Type
	pad := m.Packager.Fields[field].Padding.Pad
	lengthMessage := len(m.FieldAndSubFields[field].Field)
	lengthPackager := m.Packager.Fields[field].Length
	switch m.Packager.Fields[field].Padding.Position {
	case "RIGHT":
		padResult := padding.RightEncode(typePad, pad, lengthMessage, lengthPackager)
		return padResult, ""
	case "LEFT":
		padResult := padding.LeftEncode(typePad, pad, lengthMessage, lengthPackager)
		return "", padResult
	default:
		return "", ""
	}
}
