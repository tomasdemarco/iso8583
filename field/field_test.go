package field

import (
	"gitlab.com/g6604/adquirencia/desarrollo/golang_package/iso8583/encoding"
	"gitlab.com/g6604/adquirencia/desarrollo/golang_package/iso8583/packager"
	"gitlab.com/g6604/adquirencia/desarrollo/golang_package/iso8583/prefix"
	"testing"
)

var (
	FieldEncoding       = []encoding.Encoding{encoding.Bcd, encoding.Ascii, encoding.Ebcdic}
	FieldValuesEncoding = []string{"000001", "303030303031", "f0f0f0f0f0f1"}
	FieldPrefix         = []prefix.Type{prefix.Fixed, prefix.LL, prefix.LLL}
	FieldPrefixEncoding = []encoding.Encoding{encoding.Bcd, encoding.Hex, encoding.Ascii, encoding.Ebcdic}
	FieldPrefixValues   = []string{"", "06", "0006", "", "06", "0006", "", "3036", "303036", "", "f0f6", "f0f0f6"}
)

// TestUnpackField calls message.UnpackField
func TestUnpackField(t *testing.T) {
	for e, fieldEncoding := range FieldEncoding {
		for pe, prefixEncoding := range FieldPrefixEncoding {
			for p, fieldPrefix := range FieldPrefix {
				expectedResult := "000001"
				data := FieldPrefixValues[p+(pe*3)] + FieldValuesEncoding[e]

				fieldsPackager := packager.Field{}
				fieldsPackager.Length = 6
				fieldsPackager.Encoding = fieldEncoding
				fieldsPackager.Prefix.Type = fieldPrefix
				fieldsPackager.Prefix.Encoding = prefixEncoding

				fields := make(map[string]packager.Field)
				fields["011"] = fieldsPackager

				pkg := packager.Packager{}
				pkg.Fields = fields

				result, _, err := Unpack(fieldsPackager, data, 0, "011")
				if err != nil {
					t.Fatalf(`UnpackField(%s) Encoding=%s - Prefix=%s - PrefixEncoding=%s - Error %s`, data, fieldEncoding.String(), fieldPrefix.String(), prefixEncoding.String(), err.Error())
				}

				if *result != expectedResult {
					t.Fatalf(`UnpackField(%s) Encoding=%s - Prefix=%s - PrefixEncoding=%s - Result "%s" does not match "%s"`, data, fieldEncoding.String(), fieldPrefix.String(), prefixEncoding.String(), *result, expectedResult)
				}
				t.Logf(`UnpackField=%-28s Encoding=%-6s - Prefix=%-5s - PrefixEncoding=%-6s - Result "%-6s" match "%-6s"`, data, fieldEncoding.String(), fieldPrefix.String(), prefixEncoding.String(), *result, expectedResult)
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
				expectedResult := FieldPrefixValues[p+(pe*3)] + FieldValuesEncoding[e]

				fieldsPackager := packager.Field{}
				fieldsPackager.Length = 6
				fieldsPackager.Encoding = fieldEncoding
				fieldsPackager.Prefix.Type = fieldPrefix
				fieldsPackager.Prefix.Encoding = prefixEncoding

				fields := make(map[string]packager.Field)
				fields["011"] = fieldsPackager

				pkg := packager.Packager{}
				pkg.Fields = fields

				result := Pack(fieldsPackager, data)

				if result != expectedResult {
					t.Fatalf(`PackField(%s) Encoding=%s - Prefix=%s - PrefixEncoding=%s - Result "%s" does not match "%s"`, data, fieldEncoding.String(), fieldPrefix.String(), prefixEncoding.String(), result, expectedResult)
				}
				t.Logf(`PackField=%-11s Encoding=%-6s - Prefix=%-5s - PrefixEncoding=%-6s - Result "%-28s" match "%-6s"`, data, fieldEncoding.String(), fieldPrefix.String(), prefixEncoding.String(), result, expectedResult)
			}
		}
	}
}
