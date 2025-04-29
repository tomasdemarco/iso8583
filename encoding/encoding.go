package encoding

//
//func Unpack(encoding Encoding, value []byte) (string, error) {
//	switch encoding {
//	case encoding.Ascii:
//		return ascii.Decode(value)
//	case encoding.Ebcdic:
//		return EbcdicDecode(value), nil
//	case encoding.Bcd:
//		return Decode(value)
//	default:
//		return fmt.Sprintf("%x", value), nil
//	}
//}
//
//func Pack(encoding Encoding, value string) ([]byte, error) {
//	switch encoding {
//	case encoding.Ascii:
//		return ascii.Encode(value), nil
//	case encoding.Ebcdic:
//		return EbcdicEncode(value), nil
//	case encoding.Bcd:
//		return Encode(value)
//	default:
//		return utils.Hex2Byte(value), nil
//	}
//}
//
//func PackSubField(encoding Encoding, value string) []byte {
//	switch encoding {
//	case encoding.Ascii:
//		result := ascii.Encode(value)
//		return result
//	case encoding.Ebcdic:
//		return EbcdicEncode(value)
//	default:
//		return []byte(value)
//	}
//}
