package message

import (
	"fmt"
	"github.com/tomasdemarco/iso8583/encoding"
	"strconv"
)

func (m *Message) UnpackPrefix(field string, messageRaw string, position int, prefixLength int) (int, int) {
	switch m.Packager.Fields[field].PrefixEncoding {
	case "ASCII":
		prefixLengthAux := prefixLength
		if prefixLengthAux == 4 {
			prefixLengthAux--
		}
		lengthString, err := encoding.AsciiDecode(messageRaw[position : position+(prefixLengthAux*2)])
		length, err := strconv.Atoi(lengthString)
		if err != nil {
			fmt.Println(err)
		}
		return length, prefixLengthAux * 2
	case "HEX":
		lengthString, _ := encoding.HexDecode(messageRaw[position : position+prefixLength])
		length, err := strconv.Atoi(lengthString)
		if err != nil {
			fmt.Println(err)
		}
		return length, prefixLength
	case "EBCDIC":
		prefixLengthAux := prefixLength
		if prefixLengthAux == 4 {
			prefixLengthAux--
		}
		lengthString, _ := encoding.EbcdicDecode(messageRaw[position : position+(prefixLengthAux*2)])
		length, err := strconv.Atoi(lengthString)
		if err != nil {
			fmt.Println(err)
		}
		return length, prefixLengthAux * 2
	default:
		length, err := strconv.Atoi(messageRaw[position : position+prefixLength])
		if err != nil {
			fmt.Println(err)
		}
		return length, prefixLength
	}
}

func (m *Message) PackPrefix(field string, value int, prefixLength int) (prefix string) {
	switch m.Packager.Fields[field].PrefixEncoding {
	case "ASCII":
		if prefixLength == 4 {
			prefix = encoding.AsciiEncode(fmt.Sprintf("%03d", value))
		} else {
			prefix = encoding.AsciiEncode(fmt.Sprintf("%02d", value))
		}
		return prefix
	case "HEX":
		if prefixLength == 4 {
			prefix, _ = encoding.HexEncode(fmt.Sprintf("%d", value))
			prefix = fmt.Sprintf("%04s", prefix)
		} else {
			prefix, _ = encoding.HexEncode(fmt.Sprintf("%d", value))
			prefix = fmt.Sprintf("%02s", prefix)
		}
		return prefix
	case "EBCDIC":
		if prefixLength == 4 {
			prefix, _ = encoding.EbcdicEncode(fmt.Sprintf("%03d", value))
		} else {
			prefix, _ = encoding.EbcdicEncode(fmt.Sprintf("%02d", value))
		}
		return prefix
	default:
		if prefixLength == 4 {
			return fmt.Sprintf("%04d", value)
		} else {
			return fmt.Sprintf("%02d", value)
		}
	}
}
