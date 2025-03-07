package prefix

import (
	"errors"
	"fmt"
	enc "github.com/tomasdemarco/iso8583/encoding"
	"github.com/tomasdemarco/iso8583/utils"
	"strconv"
)

type Prefix struct {
	Type     Type         `json:"type"`
	Encoding enc.Encoding `json:"encoding"`
}

func Unpack(prefix Prefix, messageRaw []byte) (int, int, error) {
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
		lengthString := enc.EbcdicDecode(messageRaw[:prefixLength])

		length, err := strconv.Atoi(lengthString)
		if err != nil {
			return 0, 0, err
		}

		return length, prefixLength, nil
	default:
		length, err := strconv.Atoi(fmt.Sprintf("%x", messageRaw[:prefixLength]))
		if err != nil {
			fmt.Println(err)
		}

		return length, prefixLength, nil
	}
}

func Pack(prefix Prefix, value int) ([]byte, error) {
	if prefix.Type == Fixed {
		return nil, nil
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
		valueEnc, err := enc.HexEncode(fmt.Sprintf("%d", value))
		if err != nil {
			return nil, err
		}
		if prefix.Type == LLL || prefix.Type == LLLL {
			return utils.Hex2Byte(fmt.Sprintf("%04s", valueEnc)), nil
		} else {
			return utils.Hex2Byte(fmt.Sprintf("%02s", valueEnc)), nil
		}
	case enc.Ebcdic:
		if prefix.Type == LLLL {
			valueEnc := enc.EbcdicEncode(fmt.Sprintf("%04d", value))

			return valueEnc, nil
		} else if prefix.Type == LLL {
			valueEnc := enc.EbcdicEncode(fmt.Sprintf("%03d", value))

			return valueEnc, nil
		} else {
			valueEnc := enc.EbcdicEncode(fmt.Sprintf("%02d", value))

			return valueEnc, nil
		}
	default:
		if prefix.Type == LLL || prefix.Type == LLLL {
			return utils.Hex2Byte(fmt.Sprintf("%04d", value)), nil
		} else {
			return utils.Hex2Byte(fmt.Sprintf("%02d", value)), nil
		}
	}
}

func GetPrefixLen(prefixType Type, prefixEncoding enc.Encoding) int {
	var length int
	switch prefixType {
	case LL:
		if prefixEncoding == enc.Ascii || prefixEncoding == enc.Ebcdic {
			length = 2
		} else {
			length = 1
		}
	case LLL:
		if prefixEncoding == enc.Ascii || prefixEncoding == enc.Ebcdic {
			length = 3
		} else {
			length = 2
		}
	case LLLL:
		if prefixEncoding == enc.Ascii || prefixEncoding == enc.Ebcdic {
			length = 4
		} else {
			length = 2
		}
	default:
		length = 1
	}

	return length
}
