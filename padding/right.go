package padding

import (
	"github.com/tomasdemarco/iso8583/utils"
	"strings"
)

func RightDecode(paddingType Type, lengthField int) int {
	switch paddingType {
	case Parity:
		if lengthField%2 == 0 {
			return 0
		}
		return 1
	default:
		return 0
	}
}

func RightEncode(paddingType Type, char utils.ByteFromString, lengthPackager int, lengthValue int) string {
	switch paddingType {
	case Fill:
		return strings.Repeat(string(char), lengthPackager-lengthValue)
	case Parity:
		if lengthValue%2 != 0 {
			return string(char)
		}
		return ""
	default:
		return ""
	}
}
