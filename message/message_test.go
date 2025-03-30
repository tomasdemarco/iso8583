package message

import (
	"bytes"
	"github.com/tomasdemarco/iso8583/encoding"
	"github.com/tomasdemarco/iso8583/packager"
	"github.com/tomasdemarco/iso8583/packager/field"
	"github.com/tomasdemarco/iso8583/prefix"
	"slices"
	"testing"
)

var (
	FieldEncoding       = []encoding.Encoding{encoding.Bcd, encoding.Ascii, encoding.Ebcdic}
	FieldValuesEncoding = [][]byte{{0x00, 0x00, 0x01}, {0x30, 0x30, 0x30, 0x30, 0x30, 0x31}, {0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF1}}
	FieldPrefix         = []prefix.Type{prefix.Fixed, prefix.LL, prefix.LLL}
	FieldPrefixEncoding = []encoding.Encoding{encoding.Bcd, encoding.Hex, encoding.Ascii, encoding.Ebcdic}
	FieldPrefixValues   = [][]byte{{}, {0x06}, {0x00, 0x06}, {}, {0x06}, {0x00, 0x06}, {}, {0x30, 0x36}, {0x30, 0x30, 0x36}, {}, {0xF0, 0xF6}, {0xF0, 0xF0, 0xF6}}
)

// TestUnpack calls message.Unpack
func TestUnpack(t *testing.T) {
	for e, enc := range FieldEncoding {
		for pe, prefixEncoding := range FieldPrefixEncoding {
			for p, prefix2 := range FieldPrefix {
				expectedResult := "000001"
				expectedBitmap := []string{"011"}

				buf := new(bytes.Buffer)
				buf.Write([]byte{0x02, 0x00, 0x00, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
				buf.Write(FieldPrefixValues[p+(pe*3)])
				buf.Write(FieldValuesEncoding[e])
				data := buf.Bytes()

				message := Message{}
				fieldsPackager := field.Field{}
				fieldsPackager.Length = 4
				fieldsPackager.Encoding = encoding.Bcd

				fields := make(map[string]field.Field)
				fields["000"] = fieldsPackager

				fieldsPackager = field.Field{}
				fieldsPackager.Length = 8
				fieldsPackager.Encoding = encoding.Bcd

				fields["001"] = fieldsPackager

				fieldsPackager = field.Field{}
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
					t.Fatalf(`Unpack(%x) - Encoding=%s - Prefix=%s - PrefixEncoding=%s - Error %s`, data, enc.String(), prefix2.String(), prefixEncoding.String(), err.Error())
				}

				if len(message.Bitmap) != len(expectedBitmap) {
					t.Fatalf(`Unpack(%x) - Encoding=%s - Prefix=%s - PrefixEncoding=%s - Length bitmap is different - Result "%s" / Expected "%s"`, data, enc.String(), prefix2.String(), prefixEncoding.String(), message.Bitmap, expectedBitmap)
				}

				for !slices.Equal(message.Bitmap, expectedBitmap) {
					t.Fatalf(`Unpack(%x) - Encoding=%s - Prefix=%s - PrefixEncoding=%s - Result bitmap "%s" does not match "%s"`, data, enc.String(), prefix2.String(), prefixEncoding.String(), message.Bitmap, expectedBitmap)
				}

				result, err := message.GetField("011")
				if err != nil {
					t.Fatalf(`Unpack(%x) Encoding=%s - Prefix=%s - PrefixEncoding=%s - Error %s`, data, enc.String(), prefix2.String(), prefixEncoding.String(), err.Error())
				}

				if result != expectedResult {
					t.Fatalf(`Unpack(%x) Encoding=%s - Prefix=%s - PrefixEncoding=%s - Result "%s" does not match "%s"`, data, enc.String(), prefix2.String(), prefixEncoding.String(), result, expectedResult)
				}
				t.Logf(`Unpack=%x Encoding=%s - Prefix=%s - PrefixEncoding=%s - Result "%s" match "%s"`, data, enc.String(), prefix2.String(), prefixEncoding.String(), result, expectedResult)
			}
		}
	}
}
