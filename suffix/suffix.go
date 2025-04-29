package suffix

import (
	enc "github.com/tomasdemarco/iso8583/encoding"
)

type Suffix struct {
	Type     Type         `json:"type"`
	Encoding enc.Encoding `json:"encoding"`
}

//func Unpack(suffix Suffix, messageRaw []byte) (int, int, error) {
//	if suffix.Type == Fixed {
//		return 0, 0, nil
//	}
//
//	suffixLength := GetSuffixLen(suffix.Type, suffix.Encoding)
//
//	if len(messageRaw) < suffixLength {
//		return 0, 0, errors.New("index out of range while trying to unpack suffix")
//	}
//
//	switch suffix.Encoding {
//	case enc.Ascii:
//		lengthString, err := ascii.Decode(messageRaw[:suffixLength])
//		length, err := strconv.Atoi(lengthString)
//		if err != nil {
//			return 0, 0, err
//		}
//
//		return length, suffixLength, nil
//	case enc.Hex:
//		lengthString, _ := enc.HexDecode(messageRaw[:suffixLength])
//		length, err := strconv.Atoi(lengthString)
//		if err != nil {
//			fmt.Println(err)
//		}
//
//		return length, suffixLength, nil
//	case enc.Ebcdic:
//		lengthString := enc.EbcdicDecode(messageRaw[:suffixLength])
//
//		length, err := strconv.Atoi(lengthString)
//		if err != nil {
//			return 0, 0, err
//		}
//
//		return length, suffixLength, nil
//	default:
//		length, err := strconv.Atoi(fmt.Sprintf("%x", messageRaw[:suffixLength]))
//		if err != nil {
//			fmt.Println(err)
//		}
//
//		return length, suffixLength, nil
//	}
//}
//
//func Pack(suffix Suffix, value int) ([]byte, error) {
//	if suffix.Type == Fixed {
//		return nil, nil
//	}
//
//	switch suffix.Encoding {
//	case enc.Ascii:
//		if suffix.Type == LLLL {
//			return ascii.Encode(fmt.Sprintf("%04d", value)), nil
//		}
//		if suffix.Type == LLL {
//			return ascii.Encode(fmt.Sprintf("%03d", value)), nil
//		} else {
//			return ascii.Encode(fmt.Sprintf("%02d", value)), nil
//		}
//	case enc.Hex:
//		valueEnc, err := enc.HexEncode(fmt.Sprintf("%d", value))
//		if err != nil {
//			return nil, err
//		}
//		if suffix.Type == LLL || suffix.Type == LLLL {
//			return utils.Hex2Byte(fmt.Sprintf("%04s", valueEnc)), nil
//		} else {
//			return utils.Hex2Byte(fmt.Sprintf("%02s", valueEnc)), nil
//		}
//	case enc.Ebcdic:
//		if suffix.Type == LLLL {
//			valueEnc := enc.EbcdicEncode(fmt.Sprintf("%04d", value))
//
//			return valueEnc, nil
//		} else if suffix.Type == LLL {
//			valueEnc := enc.EbcdicEncode(fmt.Sprintf("%03d", value))
//
//			return valueEnc, nil
//		} else {
//			valueEnc := enc.EbcdicEncode(fmt.Sprintf("%02d", value))
//
//			return valueEnc, nil
//		}
//	default:
//		if suffix.Type == LLL || suffix.Type == LLLL {
//			return utils.Hex2Byte(fmt.Sprintf("%04d", value)), nil
//		} else {
//			return utils.Hex2Byte(fmt.Sprintf("%02d", value)), nil
//		}
//	}
//}
//
//func GetSuffixLen(suffixType Type, suffixEncoding enc.Encoding) int {
//	var length int
//	switch suffixType {
//	case LL:
//		if suffixEncoding == enc.Ascii || suffixEncoding == enc.Ebcdic {
//			length = 2
//		} else {
//			length = 1
//		}
//	case LLL:
//		if suffixEncoding == enc.Ascii || suffixEncoding == enc.Ebcdic {
//			length = 3
//		} else {
//			length = 2
//		}
//	case LLLL:
//		if suffixEncoding == enc.Ascii || suffixEncoding == enc.Ebcdic {
//			length = 4
//		} else {
//			length = 2
//		}
//	default:
//		length = 1
//	}
//
//	return length
//}
