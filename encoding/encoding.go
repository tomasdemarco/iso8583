package encoding

import (
	"fmt"
	"github.com/tomasdemarco/iso8583/utils"
)

func Unpack(encoding Encoding, value []byte) (string, error) {
	switch encoding {
	case Ascii:
		result, err := AsciiDecode(value)
		return result, err
	case Ebcdic:
		result := EbcdicDecode(value)
		return result, nil
	default:
		return fmt.Sprintf("%x", value), nil
	}
}

func Pack(encoding Encoding, value string) []byte {
	switch encoding {
	case Ascii:
		result := AsciiEncode(value)
		return result
	case Ebcdic:
		return EbcdicEncode(value)
	default:
		return utils.Hex2Byte(value)
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
