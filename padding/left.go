package padding

import (
	"strings"
)

func LeftDecode(paddingType string) int {
	switch paddingType {
	case "PARITY":
		return 1
	default:
		return 0
	}

}

func LeftEncode(paddingType string, pad string, lengthMessage int, lengthPackager int) string {
	switch paddingType {
	case "PARITY":
		if lengthMessage%2 != 0 {
			return pad
		}
		return ""
	default:
		return strings.Repeat(pad, lengthPackager-lengthMessage)
	}

}
