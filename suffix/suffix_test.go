package suffix

import (
	"bytes"
	"fmt"
	"github.com/tomasdemarco/iso8583/encoding"
	"testing"
)

var (
	Suffixes                 = []Type{LL, LLL}
	SuffixEncoding           = []encoding.Encoding{encoding.Bcd, encoding.Hex, encoding.Ascii, encoding.Ebcdic}
	SuffixValues             = [][]byte{{0x06}, {0x06}, {0x30, 0x36}, {0xF0, 0xF6}, {0x00, 0x06}, {0x00, 0x06}, {0x30, 0x30, 0x36}, {0xF0, 0xF0, 0xF6}}
	ResultSuffixValues       = []int{6, 6, 6, 6, 6, 6, 6, 6}
	ResultSuffixLengthValues = []int{1, 1, 2, 2, 2, 2, 3, 3}
)

// TestUnpackSuffix calls message.UnpackSuffix
func TestUnpackSuffix(t *testing.T) {
	for p, suffixType := range Suffixes {
		for pe, suffixEncoding := range SuffixEncoding {
			data := SuffixValues[pe+(p*4)]

			var suffix Suffix
			suffix.Type = suffixType
			suffix.Encoding = suffixEncoding

			resultSuffix, resultSuffixLength, err := Unpack(suffix, data)
			fmt.Println(resultSuffix, resultSuffixLength)
			if err != nil {
				t.Fatalf(`UnpackSuffix(%x) - Suffix=%s - SuffixEncoding=%s - Error %s`, data, suffixType.String(), suffixEncoding.String(), err.Error())
			}

			if resultSuffix != ResultSuffixValues[pe+(p*4)] {
				t.Fatalf(`UnpackSuffix(%x) - Suffix=%s - SuffixEncoding=%s - Result %d does not match %d`, data, suffixType.String(), suffixEncoding.String(), resultSuffix, ResultSuffixValues[pe+(p*4)])
			}
			t.Logf(`UnpackSuffix(%x) - Suffix=%s - SuffixEncoding=%s - Result %d match %d`, data, suffixType.String(), suffixEncoding.String(), resultSuffix, ResultSuffixValues[pe+(p*4)])

			if resultSuffixLength != ResultSuffixLengthValues[pe+(p*4)] {
				t.Fatalf(`UnpackSuffix(%x) - Suffix=%s - SuffixEncoding=%s - Result %d does not match %d`, data, suffixType.String(), suffixEncoding.String(), resultSuffixLength, ResultSuffixLengthValues[pe+(p*4)])
			}
			t.Logf(`UnpackSuffix(%x) - Suffix=%s - SuffixEncoding=%s - Result %d match %d`, data, suffixType.String(), suffixEncoding.String(), resultSuffixLength, ResultSuffixLengthValues[pe+(p*4)])
		}
	}
}

// TestPackSuffix calls encoding.PackSuffix
func TestPackSuffix(t *testing.T) {
	for p, suffixType := range Suffixes {
		for pe, suffixEncoding := range SuffixEncoding {
			data := "000001"
			expectedResult := SuffixValues[pe+(p*4)]

			var suffix Suffix
			suffix.Type = suffixType
			suffix.Encoding = suffixEncoding

			result, err := Pack(suffix, len(data))
			if err != nil {
				t.Fatalf(`PackSuffix(%x) - Suffix=%s - SuffixEncoding=%s - Error %s`, data, suffixType.String(), suffixEncoding.String(), err.Error())
			}

			if !bytes.Equal(result, expectedResult) {
				t.Fatalf(`Suffix(%s) - Suffix=%s - SuffixEncoding=%s - Result "%x" does not match "%x"`, data, suffixType.String(), suffixEncoding.String(), result, expectedResult)
			}
			t.Logf(`Suffix(%s) - Suffix=%s - SuffixEncoding=%s - Result "%x" match "%x"`, data, suffixType.String(), suffixEncoding.String(), result, expectedResult)
		}
	}
}
