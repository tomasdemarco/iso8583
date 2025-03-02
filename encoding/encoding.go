package encoding

func Unpack(encoding Encoding, value string) (string, error) {
	switch encoding {
	case Ascii:
		result, err := AsciiDecode(value)
		return result, err
	case Ebcdic:
		result, err := EbcdicDecode(value)
		return result, err
	default:
		return value, nil
	}
}

func Pack(encoding Encoding, value string) string {
	switch encoding {
	case Ascii:
		result := AsciiEncode(value)
		return result
	case Ebcdic:
		result, _ := EbcdicEncode(value)
		return result
	default:
		return value
	}
}

func PackSubField(encoding Encoding, value string) string {
	switch encoding {
	case Ascii:
		result := AsciiEncode(value)
		return result
	case Ebcdic:
		result, _ := EbcdicEncode(value)
		return result
	default:
		return value
	}
}
