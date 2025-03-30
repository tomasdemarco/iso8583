package encoding

import (
	"bytes"
)

//func BcdDecode(src []byte) (string, error) {
//	var dst bytes.Buffer
//	for _, b := range src {
//		high := b >> 4
//		low := b & 0x0F
//
//		if high > 9 || low > 9 {
//			return "", fmt.Errorf("invalid BCD byte: %x", b)
//		}
//
//		dst.WriteString(strconv.Itoa(int(high)))
//		dst.WriteString(strconv.Itoa(int(low)))
//	}
//	return dst.String(), nil
//}

//func BcdEncode(src string) ([]byte, error) {
//	if len(src)%2 != 0 {
//		return nil, fmt.Errorf("BCD string length must be even")
//	}
//
//	var result bytes.Buffer
//	for i := 0; i < len(src); i += 2 {
//		high, err := strconv.ParseUint(string(src[i]), 10, 4)
//		if err != nil {
//			return nil, fmt.Errorf("BCD string invalid digit: %s", string(src[i]))
//		}
//
//		low, err := strconv.ParseUint(string(src[i+1]), 10, 4)
//		if err != nil {
//			return nil, fmt.Errorf("BCD string invalid digit: %s", string(src[i+1]))
//		}
//
//		result.WriteByte(byte(high<<4 | low))
//	}
//
//	return result.Bytes(), nil
//}

// str2bcd convierte una cadena a BCD.
// func str2bcd(s string, padLeft bool) []byte {
func BcdEncode(s string) ([]byte, error) {
	start := 0
	d := make([]byte, (len(s)+1)/2)

	//if len%2 == 1 && padLeft {
	//	start = 1
	//}

	for i := start; i < len(s)+start; i++ {
		n := i / 2
		digit := s[i-start] - '0'
		if i%2 == 1 {
			d[n] |= digit
		} else {
			d[n] |= digit << 4
		}
	}
	return d, nil
}

// bcd2str convierte BCD a una cadena.
// func bcd2str(b []byte, offset, len int, padLeft bool) string {
func BcdDecode(b []byte) (string, error) {
	start := 0
	var d bytes.Buffer

	//if len%2 == 1 && padLeft {
	//	start = 1
	//}

	for i := start; i < len(b)*2+start; i++ {
		shift := 0
		if i%2 == 1 {
			shift = 0
		} else {
			shift = 4
		}

		c := (b[i/2] >> shift) & 0xF
		var char rune
		if c < 10 {
			char = rune(c + '0')
		} else {
			char = rune(c - 10 + 'A')
		}

		if char == 'D' {
			char = '='
		}
		d.WriteRune(char)
	}
	return d.String(), nil
}
