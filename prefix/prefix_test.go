package prefix

import (
	"bytes"
	"fmt"
	"github.com/tomasdemarco/iso8583/encoding"
	"testing"
)

var (
	Prefixes                 = []Type{LL, LLL}
	PrefixEncoding           = []encoding.Encoding{encoding.Bcd, encoding.Hex, encoding.Ascii, encoding.Ebcdic}
	PrefixValues             = [][]byte{{0x06}, {0x06}, {0x30, 0x36}, {0xF0, 0xF6}, {0x00, 0x06}, {0x00, 0x06}, {0x30, 0x30, 0x36}, {0xF0, 0xF0, 0xF6}}
	ResultPrefixValues       = []int{6, 6, 6, 6, 6, 6, 6, 6}
	ResultPrefixLengthValues = []int{1, 1, 2, 2, 2, 2, 3, 3}
)

// TestUnpackPrefix calls message.UnpackPrefix
func TestUnpackPrefix(t *testing.T) {
	for p, prefixType := range Prefixes {
		for pe, prefixEncoding := range PrefixEncoding {
			data := PrefixValues[pe+(p*4)]

			var prefix Prefix
			prefix.Type = prefixType
			prefix.Encoding = prefixEncoding

			resultPrefix, resultPrefixLength, err := Unpack(prefix, data)
			fmt.Println(resultPrefix, resultPrefixLength)
			if err != nil {
				t.Fatalf(`UnpackPrefix(%x) - Prefix=%s - PrefixEncoding=%s - Error %s`, data, prefixType.String(), prefixEncoding.String(), err.Error())
			}

			if resultPrefix != ResultPrefixValues[pe+(p*4)] {
				t.Fatalf(`UnpackPrefix(%x) - Prefix=%s - PrefixEncoding=%s - Result %d does not match %d`, data, prefixType.String(), prefixEncoding.String(), resultPrefix, ResultPrefixValues[pe+(p*4)])
			}
			t.Logf(`UnpackPrefix(%x) - Prefix=%s - PrefixEncoding=%s - Result %d match %d`, data, prefixType.String(), prefixEncoding.String(), resultPrefix, ResultPrefixValues[pe+(p*4)])

			if resultPrefixLength != ResultPrefixLengthValues[pe+(p*4)] {
				t.Fatalf(`UnpackPrefix(%x) - Prefix=%s - PrefixEncoding=%s - Result %d does not match %d`, data, prefixType.String(), prefixEncoding.String(), resultPrefixLength, ResultPrefixLengthValues[pe+(p*4)])
			}
			t.Logf(`UnpackPrefix(%x) - Prefix=%s - PrefixEncoding=%s - Result %d match %d`, data, prefixType.String(), prefixEncoding.String(), resultPrefixLength, ResultPrefixLengthValues[pe+(p*4)])
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

			result, err := Pack(prefix, len(data))
			if err != nil {
				t.Fatalf(`PackPrefix(%x) - Prefix=%s - PrefixEncoding=%s - Error %s`, data, prefixType.String(), prefixEncoding.String(), err.Error())
			}

			if !bytes.Equal(result, expectedResult) {
				t.Fatalf(`Prefix(%s) - Prefix=%s - PrefixEncoding=%s - Result "%x" does not match "%x"`, data, prefixType.String(), prefixEncoding.String(), result, expectedResult)
			}
			t.Logf(`Prefix(%s) - Prefix=%s - PrefixEncoding=%s - Result "%x" match "%x"`, data, prefixType.String(), prefixEncoding.String(), result, expectedResult)
		}
	}
}
