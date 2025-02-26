package prefix

import (
	"fmt"
	enc "github.com/tomasdemarco/iso8583/encoding"
	"strconv"
)

type Prefix struct {
	Type     Type         `json:"type"`
	Encoding enc.Encoding `json:"encoding"`
}

func Unpack(prefix Prefix, messageRaw string) (int, int, error) {

	prefixLength := GetPrefixLen(prefix.Type, prefix.Encoding)

	switch prefix.Encoding {
	case enc.Ascii:
		lengthString, err := enc.AsciiDecode(messageRaw[:prefixLength])
		length, err := strconv.Atoi(lengthString)
		if err != nil {
			return 0, 0, err
		}

		return length, prefixLength, nil
	case enc.Hex:
		lengthString, _ := enc.HexDecode(messageRaw[:prefixLength])
		length, err := strconv.Atoi(lengthString)
		if err != nil {
			fmt.Println(err)
		}

		return length, prefixLength, nil
	case enc.Ebcdic:
		lengthString, err := enc.EbcdicDecode(messageRaw[:prefixLength])
		if err != nil {
			return 0, 0, err
		}

		length, err := strconv.Atoi(lengthString)
		if err != nil {
			return 0, 0, err
		}

		return length, prefixLength, nil
	default:
		length, err := strconv.Atoi(messageRaw[:prefixLength])
		if err != nil {
			fmt.Println(err)
		}

		return length, prefixLength, nil
	}
}

func Pack(prefix Prefix, value int) (string, error) {
	switch prefix.Encoding {
	case enc.Ascii:
		if prefix.Type == LLL {
			return enc.AsciiEncode(fmt.Sprintf("%03d", value)), nil
		} else {
			return enc.AsciiEncode(fmt.Sprintf("%02d", value)), nil
		}
	case enc.Hex:
		if prefix.Type == LLL {
			prefix, err := enc.HexEncode(fmt.Sprintf("%d", value))
			if err != nil {
				return "", err
			}

			return fmt.Sprintf("%04s", prefix), nil
		} else {
			prefix, err := enc.HexEncode(fmt.Sprintf("%d", value))
			if err != nil {
				return "", err
			}

			return fmt.Sprintf("%02s", prefix), nil
		}
	case enc.Ebcdic:
		if prefix.Type == LLL {
			prefix, err := enc.EbcdicEncode(fmt.Sprintf("%03d", value))
			if err != nil {
				return "", err
			}

			return prefix, nil
		} else {
			prefix, err := enc.EbcdicEncode(fmt.Sprintf("%02d", value))
			if err != nil {
				return "", err
			}

			return prefix, nil
		}
	default:
		if prefix.Type == LLL {
			return fmt.Sprintf("%04d", value), nil
		} else {
			return fmt.Sprintf("%02d", value), nil
		}
	}
}

func GetPrefixLen(prefixType Type, prefixEncoding enc.Encoding) int {
	var length int
	switch prefixType {
	case LL:
		length = 2
	case LLL:
		length = 4
	default:
		length = 2
	}

	if prefixEncoding == enc.Ascii ||
		prefixEncoding == enc.Ebcdic {
		if prefixType == LLL {
			length--
		}

		length = length * 2
	}

	return length
}
