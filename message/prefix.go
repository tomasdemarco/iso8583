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
			panic(err)
		}
		return length * 2, prefixLengthAux * 2
	case "HEX":
		if m.Packager.Fields[field].Type == "STRING" {
			lengthString := encoding.HexDecode(messageRaw[position : position+prefixLength])
			length, err := strconv.Atoi(lengthString)
			if err != nil {
				panic(err)
			}
			return length * 2, prefixLength
		} else {
			lengthString := encoding.HexDecode(messageRaw[position : position+prefixLength])
			length, err := strconv.Atoi(lengthString)
			if err != nil {
				panic(err)
			}
			return length, prefixLength
		}
	case "EBCDIC":
		if m.Packager.Fields[field].Type == "STRING" {
			prefixLengthAux := prefixLength
			if prefixLengthAux == 4 {
				prefixLengthAux--
			}
			lengthString, _ := encoding.EbcdicDecode(messageRaw[position : position+(prefixLengthAux*2)])
			length, err := strconv.Atoi(lengthString)
			if err != nil {
				panic(err)
			}
			return length * 2, prefixLengthAux * 2
		} else {
			prefixLengthAux := prefixLength
			if prefixLengthAux == 4 {
				prefixLengthAux--
			}
			lengthString, _ := encoding.EbcdicDecode(messageRaw[position : position+(prefixLengthAux*2)])
			length, err := strconv.Atoi(lengthString)
			if err != nil {
				panic(err)
			}
			return length, prefixLengthAux * 2
		}
	default:
		length, err := strconv.Atoi(messageRaw[position : position+prefixLength])
		if err != nil {
			panic(err)
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
		if m.Packager.Fields[field].Type == "STRING" {
			if prefixLength == 4 {
				prefix = encoding.HexEncode(fmt.Sprintf("%04d", value/2))
			} else {
				prefix = encoding.HexEncode(fmt.Sprintf("%02d", value/2))
			}
			return prefix
		} else {
			if prefixLength == 4 {
				prefix = encoding.HexEncode(fmt.Sprintf("%04d", value))
			} else {
				prefix = encoding.HexEncode(fmt.Sprintf("%02d", value))
			}
			return prefix
		}
	case "EBCDIC":
		if m.Packager.Fields[field].Type == "STRING" {
			if prefixLength == 4 {
				prefix, _ = encoding.EbcdicEncode(fmt.Sprintf("%03d", value/2))
			} else {
				prefix, _ = encoding.EbcdicEncode(fmt.Sprintf("%02d", value/2))
			}
			return prefix
		} else {
			if prefixLength == 4 {
				prefix, _ = encoding.EbcdicEncode(fmt.Sprintf("%03d", value))
			} else {
				prefix, _ = encoding.EbcdicEncode(fmt.Sprintf("%02d", value))
			}
			return prefix
		}
	default:
		return fmt.Sprintf("%d", value)
	}
}
