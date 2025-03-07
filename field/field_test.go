package field

import (
	"bytes"
	"github.com/tomasdemarco/iso8583/encoding"
	"github.com/tomasdemarco/iso8583/packager"
	"github.com/tomasdemarco/iso8583/prefix"
	"testing"
)

var (
	FieldEncoding       = []encoding.Encoding{encoding.Bcd, encoding.Ascii, encoding.Ebcdic}
	FieldValuesEncoding = [][]byte{{0x00, 0x00, 0x01}, {0x30, 0x30, 0x30, 0x30, 0x30, 0x31}, {0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF1}}
	FieldPrefix         = []prefix.Type{prefix.Fixed, prefix.LL, prefix.LLL}
	FieldPrefixEncoding = []encoding.Encoding{encoding.Bcd, encoding.Hex, encoding.Ascii, encoding.Ebcdic}
	FieldPrefixValues   = [][]byte{{}, {0x06}, {0x00, 0x06}, {}, {0x06}, {0x00, 0x06}, {}, {0x30, 0x36}, {0x30, 0x30, 0x36}, {}, {0xF0, 0xF6}, {0xF0, 0xF0, 0xF6}}
)

// TestUnpackField calls message.UnpackField
func TestUnpackField(t *testing.T) {
	for e, fieldEncoding := range FieldEncoding {
		for pe, prefixEncoding := range FieldPrefixEncoding {
			for p, fieldPrefix := range FieldPrefix {
				expectedResult := "000001"
				data := append(FieldPrefixValues[p+(pe*3)], FieldValuesEncoding[e]...)

				fieldsPackager := packager.Field{}
				fieldsPackager.Length = 6
				fieldsPackager.Encoding = fieldEncoding
				if fieldPrefix != prefix.Fixed {
					fieldsPackager.Prefix.Type = fieldPrefix
					fieldsPackager.Prefix.Encoding = prefixEncoding
				}

				fields := make(map[string]packager.Field)
				fields["011"] = fieldsPackager

				pkg := packager.Packager{}
				pkg.Fields = fields

				result, _, err := Unpack(fieldsPackager, data, 0, "011")
				if err != nil {
					t.Fatalf(`UnpackField(%x) Encoding=%s - Prefix=%s - PrefixEncoding=%s - Error %s`, data, fieldEncoding.String(), fieldPrefix.String(), prefixEncoding.String(), err.Error())
				}

				if *result != expectedResult {
					t.Fatalf(`UnpackField(%x) Encoding=%s - Prefix=%s - PrefixEncoding=%s - Result "%s" does not match "%s"`, data, fieldEncoding.String(), fieldPrefix.String(), prefixEncoding.String(), *result, expectedResult)
				}
				t.Logf(`UnpackField(%x) Encoding=%s - Prefix=%s - PrefixEncoding=%s - Result "%s" match "%s"`, data, fieldEncoding.String(), fieldPrefix.String(), prefixEncoding.String(), *result, expectedResult)
			}
		}
	}
}

// TestPackField calls encoding.PackField
func TestPackField(t *testing.T) {
	for e, fieldEncoding := range FieldEncoding {
		for pe, prefixEncoding := range FieldPrefixEncoding {
			for p, fieldPrefix := range FieldPrefix {
				data := "000001"
				expectedResult := append(FieldPrefixValues[p+(pe*3)], FieldValuesEncoding[e]...)

				fieldsPackager := packager.Field{}
				fieldsPackager.Length = 6
				fieldsPackager.Encoding = fieldEncoding
				if fieldPrefix != prefix.Fixed {
					fieldsPackager.Prefix.Type = fieldPrefix
					fieldsPackager.Prefix.Encoding = prefixEncoding
				}

				fields := make(map[string]packager.Field)
				fields["011"] = fieldsPackager

				pkg := packager.Packager{}
				pkg.Fields = fields

				result, err := Pack(fieldsPackager, data)
				if err != nil {
					t.Fatalf(`PackField(%s) Encoding=%s - Prefix=%s - PrefixEncoding=%s - Error %s`, data, fieldEncoding.String(), fieldPrefix.String(), prefixEncoding.String(), err.Error())
				}

				if !bytes.Equal(result, expectedResult) {
					t.Fatalf(`PackField(%s) Encoding=%s - Prefix=%s - PrefixEncoding=%s - Result "%x" does not match "%x"`, data, fieldEncoding.String(), fieldPrefix.String(), prefixEncoding.String(), result, expectedResult)
				}
				t.Logf(`PackField(%s) Encoding=%s - Prefix=%s - PrefixEncoding=%s - Result "%x" match "%x"`, data, fieldEncoding.String(), fieldPrefix.String(), prefixEncoding.String(), result, expectedResult)
			}
		}
	}
}
