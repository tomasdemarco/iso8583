package message

import (
	"github.com/tomasdemarco/iso8583/packager"
	"testing"
)

var (
	FieldEncoding       = []string{"BCD", "ASCII", "EBCDIC"}
	FieldValuesEncoding = []string{"000001", "303030303031", "f0f0f0f0f0f1"}
	FieldPrefix         = []string{"FIXED", "LL", "LLL"}
	FieldPrefixEncoding = []string{"BCD", "HEX", "ASCII", "EBCDIC"}
	FieldPrefixValues   = []string{"", "06", "0006", "", "06", "0006", "", "3036", "303036", "", "f0f6", "f0f0f6"}
)

// TestUnpackField calls message.UnpackField
func TestUnpackField(t *testing.T) {
	for e, encoding := range FieldEncoding {
		for pe, prefixEncoding := range FieldPrefixEncoding {
			for p, prefix := range FieldPrefix {
				expectedResult := "000001"
				data := FieldPrefixValues[p+(pe*3)] + FieldValuesEncoding[e]

				message := Message{}
				fieldsPackager := packager.FieldsPackager{}
				fieldsPackager.Length = 6
				fieldsPackager.Encoding = encoding
				fieldsPackager.Prefix = prefix
				fieldsPackager.PrefixEncoding = prefixEncoding

				fields := make(map[string]packager.FieldsPackager)
				fields["011"] = fieldsPackager

				pkg := packager.Packager{}
				pkg.Fields = fields
				message.Packager = &pkg

				_, err := message.UnpackField(data, 0, "011")
				if err != nil {
					t.Fatalf(`UnpackField(%s) Encoding=%s - Prefix=%s - PrefixEncoding=%s - Error %s`, data, encoding, prefix, prefixEncoding, err.Error())
				}

				result, err := message.GetField("011")
				if err != nil {
					t.Fatalf(`UnpackField(%s) Encoding=%s - Prefix=%s - PrefixEncoding=%s - Error %s`, data, encoding, prefix, prefixEncoding, err.Error())
				}

				if result != expectedResult {
					t.Fatalf(`UnpackField(%s) Encoding=%s - Prefix=%s - PrefixEncoding=%s - Result "%s" does not match "%s"`, data, encoding, prefix, prefixEncoding, result, expectedResult)
				}
				t.Logf(`UnpackField=%-28s Encoding=%-6s - Prefix=%-5s - PrefixEncoding=%-6s - Result "%-6s" match "%-6s"`, data, encoding, prefix, prefixEncoding, result, expectedResult)
			}
		}
	}
}

// TestPackField calls encoding.PackField
func TestPackField(t *testing.T) {
	for e, encoding := range FieldEncoding {
		for pe, prefixEncoding := range FieldPrefixEncoding {
			for p, prefix := range FieldPrefix {
				data := "000001"
				expectedResult := FieldPrefixValues[p+(pe*3)] + FieldValuesEncoding[e]

				message := Message{}
				fieldsPackager := packager.FieldsPackager{}
				fieldsPackager.Length = 6
				fieldsPackager.Encoding = encoding
				fieldsPackager.Prefix = prefix
				fieldsPackager.PrefixEncoding = prefixEncoding

				fields := make(map[string]packager.FieldsPackager)
				fields["011"] = fieldsPackager

				pkg := packager.Packager{}
				pkg.Fields = fields
				message.Packager = &pkg

				message.SetField("011", data)

				result := message.PackField("011")

				if result != expectedResult {
					t.Fatalf(`PackField(%s) Encoding=%s - Prefix=%s - PrefixEncoding=%s - Result "%s" does not match "%s"`, data, encoding, prefix, prefixEncoding, result, expectedResult)
				}
				t.Logf(`PackField=%-11s Encoding=%-6s - Prefix=%-5s - PrefixEncoding=%-6s - Result "%-28s" match "%-6s"`, data, encoding, prefix, prefixEncoding, result, expectedResult)
			}
		}
	}
}
