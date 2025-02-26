package message

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

// TestUnpack calls message.Unpack
func TestUnpack(t *testing.T) {
	for e, enc := range FieldEncoding {
		for pe, prefixEncoding := range FieldPrefixEncoding {
			for p, prefix2 := range FieldPrefix {
				expectedResult := "000001"
				expectedBitmap := []string{"011"}
				data := "02000020000000000000" + FieldPrefixValues[p+(pe*3)] + FieldValuesEncoding[e]

				message := Message{}
				fieldsPackager := packager.Field{}
				fieldsPackager.Length = 4
				fieldsPackager.Encoding = encoding.Bcd

				fields := make(map[string]packager.Field)
				fields["000"] = fieldsPackager

				fieldsPackager = packager.Field{}
				fieldsPackager.Length = 16
				fieldsPackager.Encoding = encoding.Bcd

				fields["001"] = fieldsPackager

				fieldsPackager = packager.Field{}
				fieldsPackager.Length = 6
				fieldsPackager.Encoding = enc
				fieldsPackager.Prefix.Type = prefix2
				fieldsPackager.Prefix.Encoding = prefixEncoding

				fields["011"] = fieldsPackager

				pkg := packager.Packager{}
				pkg.Fields = fields
				message.Packager = &pkg

				err := message.Unpack(data)
				if err != nil {
					t.Fatalf(`Unpack(%s) Encoding=%s - Prefix=%s - PrefixEncoding=%s - Error %s`, data, enc.String(), prefix2.String(), prefixEncoding.String(), err.Error())
				}

				if len(message.Bitmap) != len(expectedBitmap) {
					t.Fatalf(`Unpack(%s) - Length bitmap is different - Result "%s" / Expected "%s"`, data, message.Bitmap, expectedBitmap)
				}

				for i, v := range message.Bitmap {
					if v != expectedBitmap[i] {
						t.Fatalf(`Unpack(%s) - Result bitmap "%s" does not match "%s"`, data, message.Bitmap, expectedBitmap)

					}
				}

				result, err := message.GetField("011")
				if err != nil {
					t.Fatalf(`Unpack(%s) Encoding=%s - Prefix=%s - PrefixEncoding=%s - Error %s`, data, enc.String(), prefix2.String(), prefixEncoding.String(), err.Error())
				}

				if result != expectedResult {
					t.Fatalf(`Unpack(%s) Encoding=%s - Prefix=%s - PrefixEncoding=%s - Result "%s" does not match "%s"`, data, enc.String(), prefix2.String(), prefixEncoding.String(), result, expectedResult)
				}
				t.Logf(`Unpack=%-28s Encoding=%-6s - Prefix=%-5s - PrefixEncoding=%-6s - Result "%-6s" match "%-6s"`, data, enc.String(), prefix2.String(), prefixEncoding.String(), result, expectedResult)
			}
		}
	}
}
