package encoding

import (
	"fmt"
	"github.com/tomasdemarco/iso8583/utils"
)

func Unpack(encoding Encoding, value []byte) (string, error) {
	switch encoding {
	case Ascii:
		return AsciiDecode(value)
	case Ebcdic:
		return EbcdicDecode(value), nil
	case Bcd:
		return BcdDecode(value)
	default:
		return fmt.Sprintf("%x", value), nil
	}
}

func Pack(encoding Encoding, value string) ([]byte, error) {
	switch encoding {
	case Ascii:
		return AsciiEncode(value), nil
	case Ebcdic:
		return EbcdicEncode(value), nil
	case Bcd:
		return BcdEncode(value)
	default:
		return utils.Hex2Byte(value), nil
	}
}

func PackSubField(encoding Encoding, value string) []byte {
	switch encoding {
	case Ascii:
		result := AsciiEncode(value)
		return result
	case Ebcdic:
		return EbcdicEncode(value)
	default:
		return []byte(value)
	}
}
