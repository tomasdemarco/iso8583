package padding

import (
	"strings"
)

func LeftDecode(typePadding string) int {
	switch typePadding {
	case "PARITY":
		return 1
	default:
		return 0
	}

}

func LeftEncode(typePadding string, pad string, lengthMessage int, lengthPackager int) string {
	switch typePadding {
	case "PARITY":
		if lengthMessage%2 != 0 {
			return pad
		}
		return ""
	default:
		return strings.Repeat(pad, lengthPackager-lengthMessage)
	}

}
