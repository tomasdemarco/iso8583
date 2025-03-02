package prefix

import (
	"errors"
	"fmt"
	enc "github.com/tomasdemarco/iso8583/encoding"
	"strconv"
)

type Prefix struct {
	Type     Type         `json:"type"`
	Encoding enc.Encoding `json:"encoding"`
}

func Unpack(prefix Prefix, messageRaw string) (int, int, error) {
	if prefix.Type == Fixed {
		return 0, 0, nil
	}

	prefixLength := GetPrefixLen(prefix.Type, prefix.Encoding)

	if len(messageRaw) < prefixLength {
		return 0, 0, errors.New("index out of range while trying to unpack prefix")
	}

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
	if prefix.Type == Fixed {
		return "", nil
	}

	switch prefix.Encoding {
	case enc.Ascii:
		if prefix.Type == LLLL {
			return enc.AsciiEncode(fmt.Sprintf("%04d", value)), nil
		}
		if prefix.Type == LLL {
			return enc.AsciiEncode(fmt.Sprintf("%03d", value)), nil
		} else {
			return enc.AsciiEncode(fmt.Sprintf("%02d", value)), nil
		}
	case enc.Hex:
		if prefix.Type == LLL || prefix.Type == LLLL {
			valueEnc, err := enc.HexEncode(fmt.Sprintf("%04d", value))
			if err != nil {
				return "", err
			}

			return fmt.Sprintf("%04s", valueEnc), nil
		} else {
			valueEnc, err := enc.HexEncode(fmt.Sprintf("%02d", value))
			if err != nil {
				return "", err
			}

			return fmt.Sprintf("%02s", valueEnc), nil
		}
	case enc.Ebcdic:
		if prefix.Type == LLLL {
			valueEnc, err := enc.EbcdicEncode(fmt.Sprintf("%04d", value))
			if err != nil {
				return "", err
			}

			return valueEnc, nil
		} else if prefix.Type == LLL {
			valueEnc, err := enc.EbcdicEncode(fmt.Sprintf("%03d", value))
			if err != nil {
				return "", err
			}

			return valueEnc, nil
		} else {
			valueEnc, err := enc.EbcdicEncode(fmt.Sprintf("%02d", value))
			if err != nil {
				return "", err
			}

			return valueEnc, nil
		}
	default:
		if prefix.Type == LLL || prefix.Type == LLLL {
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
	case LLLL:
		length = 4
	default:
		length = 2
	}

	if prefixEncoding != enc.Bcd &&
		prefixEncoding != enc.Hex {
		if prefixType == LLL {
			length--
		}

		length = length * 2
	}

	return length
}
