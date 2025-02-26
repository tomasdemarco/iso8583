package prefix

import (
	"gitlab.com/g6604/adquirencia/desarrollo/golang_package/iso8583/encoding"
	"testing"
)

var (
	Prefixes                 = []Type{LL, LLL}
	PrefixEncoding           = []encoding.Encoding{encoding.Bcd, encoding.Hex, encoding.Ascii, encoding.Ebcdic}
	PrefixValues             = []string{"06", "06", "3036", "f0f6", "0006", "0006", "303036", "f0f0f6"}
	ResultPrefixValues       = []int{6, 6, 6, 6, 6, 6, 6, 6}
	ResultPrefixLengthValues = []int{2, 2, 4, 4, 4, 4, 6, 6}
)

// TestUnpackPrefix calls message.UnpackPrefix
func TestUnpackPrefix(t *testing.T) {
	for p, prefixType := range Prefixes {
		for pe, prefixEncoding := range PrefixEncoding {
			data := PrefixValues[pe+(p*4)]

			var prefix Prefix
			prefix.Type = prefixType
			prefix.Encoding = prefixEncoding

			resultPrefix, resultPrefixLength, _ := Unpack(prefix, data)

			if resultPrefix != ResultPrefixValues[pe+(p*4)] {
				t.Fatalf(`UnpackPrefix(%s) - Prefix=%s - PrefixEncoding=%s - Result %d does not match %d`, data, prefixType.String(), prefixEncoding.String(), resultPrefix, ResultPrefixValues[pe+(p*4)])
			}
			t.Logf(`UnpackPrefix=%-6s - Prefix=%s - PrefixEncoding=%-6s - Result %d match %d`, data, prefixType.String(), prefixEncoding.String(), resultPrefix, ResultPrefixValues[pe+(p*4)])

			if resultPrefixLength != ResultPrefixLengthValues[pe+(p*4)] {
				t.Fatalf(`UnpackPrefix(%s) - Prefix=%s - PrefixEncoding=%s - Result %d does not match %d`, data, prefixType.String(), prefixEncoding.String(), resultPrefixLength, ResultPrefixLengthValues[pe+(p*4)])
			}
			t.Logf(`UnpackPrefix=%-6s - Prefix=%s - PrefixEncoding=%-6s - Result %d match %d`, data, prefixType.String(), prefixEncoding.String(), resultPrefixLength, ResultPrefixLengthValues[pe+(p*4)])
		}
	}
}

// TestPackPrefix calls encoding.PackPrefix
func TestPackPrefix(t *testing.T) {
	for p, prefixType := range Prefixes {
		for pe, prefixEncoding := range PrefixEncoding {
			data := "000001"
			expectedResult := PrefixValues[pe+(p*4)]

			var prefix Prefix
			prefix.Type = prefixType
			prefix.Encoding = prefixEncoding

			result, _ := Pack(prefix, len(data))

			if result != expectedResult {
				t.Fatalf(`PackField(%s) - Prefix=%s - PrefixEncoding=%s - Result "%s" does not match "%s"`, data, prefixType.String(), prefixEncoding.String(), result, expectedResult)
			}
			t.Logf(`PackField=%-6s - Prefix=%s - PrefixEncoding=%s - Result "%s" match "%s"`, data, prefixType.String(), prefixEncoding.String(), result, expectedResult)
		}
	}
}
