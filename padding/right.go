package padding

import "strings"

func RightDecode(paddingType Type) int {
	switch paddingType {
	case Parity:
		return 1
	default:
		return 0
	}
}

func RightEncode(paddingType Type, char byte, lengthMessage int, lengthPackager int) string {
	switch paddingType {
	case Parity:
		if lengthMessage%2 != 0 {
			return string(char)
		}
		return ""
	default:
		return strings.Repeat(string(char), lengthPackager-lengthMessage)
	}
}
