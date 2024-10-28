package message

import (
	"github.com/tomasdemarco/iso8583/packager"
	"testing"
)

var (
	Prefix                   = []string{"LL", "LLL"}
	PrefixEncoding           = []string{"BCD", "HEX", "ASCII", "EBCDIC"}
	PrefixValues             = []string{"06", "06", "3036", "f0f6", "0006", "0006", "303036", "f0f0f6"}
	ResultPrefixValues       = []int{6, 6, 6, 6, 6, 6, 6, 6}
	ResultPrefixLengthValues = []int{2, 2, 4, 4, 4, 4, 6, 6}
)

// TestUnpackPrefix calls message.UnpackPrefix
func TestUnpackPrefix(t *testing.T) {
	for p, prefix := range Prefix {
		for pe, prefixEncoding := range PrefixEncoding {
			data := PrefixValues[pe+(p*4)]

			message := Message{}
			fieldsPackager := packager.FieldsPackager{}
			fieldsPackager.PrefixEncoding = prefixEncoding

			fields := make(map[string]packager.FieldsPackager)
			fields["011"] = fieldsPackager

			pkg := packager.Packager{}
			pkg.Fields = fields
			message.Packager = &pkg

			prefixLength := 2
			if prefix == "LLL" {
				prefixLength = 4
			}

			resultPrefix, resultPrefixLength := message.UnpackPrefix("011", data, 0, prefixLength)

			if resultPrefix != ResultPrefixValues[pe+(p*4)] {
				t.Fatalf(`UnpackPrefix(%s) - Prefix=%s - PrefixEncoding=%s - Result %d does not match %d`, data, prefix, prefixEncoding, resultPrefix, ResultPrefixValues[pe+(p*4)])
			}
			t.Logf(`UnpackPrefix=%-6s - Prefix=%s - PrefixEncoding=%-6s - Result %d match %d`, data, prefix, prefixEncoding, resultPrefix, ResultPrefixValues[pe+(p*4)])

			if resultPrefixLength != ResultPrefixLengthValues[pe+(p*4)] {
				t.Fatalf(`UnpackPrefix(%s) - Prefix=%s - PrefixEncoding=%s - Result %d does not match %d`, data, prefix, prefixEncoding, resultPrefixLength, ResultPrefixLengthValues[pe+(p*4)])
			}
			t.Logf(`UnpackPrefix=%-6s - Prefix=%s - PrefixEncoding=%-6s - Result %d match %d`, data, prefix, prefixEncoding, resultPrefixLength, ResultPrefixLengthValues[pe+(p*4)])
		}
	}
}

// TestPackPrefix calls encoding.PackPrefix
func TestPackPrefix(t *testing.T) {
	for p, prefix := range Prefix {
		for pe, prefixEncoding := range PrefixEncoding {
			data := "000001"
			expectedResult := PrefixValues[pe+(p*4)]

			message := Message{}
			fieldsPackager := packager.FieldsPackager{}
			fieldsPackager.PrefixEncoding = prefixEncoding

			fields := make(map[string]packager.FieldsPackager)
			fields["011"] = fieldsPackager

			pkg := packager.Packager{}
			pkg.Fields = fields
			message.Packager = &pkg

			message.SetField("011", data)

			prefixLength := 2
			if prefix == "LLL" {
				prefixLength = 4
			}

			result := message.PackPrefix("011", 6, prefixLength)

			if result != expectedResult {
				t.Fatalf(`PackField(%s) - Prefix=%s - PrefixEncoding=%s - Result "%s" does not match "%s"`, data, prefix, prefixEncoding, result, expectedResult)
			}
			t.Logf(`PackField=%-6s - Prefix=%s - PrefixEncoding=%s - Result "%s" match "%s"`, data, prefix, prefixEncoding, result, expectedResult)
		}
	}
}
