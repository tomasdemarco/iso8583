package message

import (
	"github.com/tomasdemarco/iso8583/packager"
	"testing"
)

// TestUnpack calls message.Unpack
func TestUnpack(t *testing.T) {
	for e, encoding := range FieldEncoding {
		for pe, prefixEncoding := range FieldPrefixEncoding {
			for p, prefix := range FieldPrefix {
				expectedResult := "000001"
				expectedBitmap := []string{"011"}
				data := "02000020000000000000" + FieldPrefixValues[p+(pe*3)] + FieldValuesEncoding[e]

				message := Message{}
				fieldsPackager := packager.FieldsPackager{}
				fieldsPackager.Length = 4
				fieldsPackager.Encoding = "BCD"

				fields := make(map[string]packager.FieldsPackager)
				fields["000"] = fieldsPackager

				fieldsPackager = packager.FieldsPackager{}
				fieldsPackager.Length = 16
				fieldsPackager.Encoding = "BCD"

				fields["001"] = fieldsPackager

				fieldsPackager = packager.FieldsPackager{}
				fieldsPackager.Length = 6
				fieldsPackager.Encoding = encoding
				fieldsPackager.Prefix = prefix
				fieldsPackager.PrefixEncoding = prefixEncoding

				fields["011"] = fieldsPackager

				pkg := packager.Packager{}
				pkg.Fields = fields
				message.Packager = &pkg

				err := message.Unpack(data)
				if err != nil {
					t.Fatalf(`Unpack(%s) Encoding=%s - Prefix=%s - PrefixEncoding=%s - Error %s`, data, encoding, prefix, prefixEncoding, err.Error())
				}

				if len(message.Bitmap) != len(expectedBitmap) {
					t.Fatalf(`Unpack(%s) - Length bitmap is different - Result "%s" / Expected "%s"`, data, message.Bitmap, expectedBitmap)
				}

				for i, v := range message.Bitmap {
					if v != expectedBitmap[i] {
						t.Fatalf(`Unpack(%s) - Result bitmap "%s" does not match "%s"`, data, message.Bitmap, expectedResult)

					}
				}

				result, err := message.GetField("011")
				if err != nil {
					t.Fatalf(`Unpack(%s) Encoding=%s - Prefix=%s - PrefixEncoding=%s - Error %s`, data, encoding, prefix, prefixEncoding, err.Error())
				}

				if result != expectedResult {
					t.Fatalf(`Unpack(%s) Encoding=%s - Prefix=%s - PrefixEncoding=%s - Result "%s" does not match "%s"`, data, encoding, prefix, prefixEncoding, result, expectedResult)
				}
				t.Logf(`Unpack=%-28s Encoding=%-6s - Prefix=%-5s - PrefixEncoding=%-6s - Result "%-6s" match "%-6s"`, data, encoding, prefix, prefixEncoding, result, expectedResult)
			}
		}
	}
}
