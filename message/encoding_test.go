package message

import (
	"github.com/tomasdemarco/iso8583/packager"
	"testing"
)

var (
	Encoding       = []string{"BCD", "ASCII", "EBCDIC"}
	ValuesEncoding = []string{"000001", "303030303031", "f0f0f0f0f0f1"}
)

// TestUnpackEncoding calls message.UnpackEncoding
func TestUnpackEncoding(t *testing.T) {
	for e, encoding := range Encoding {
		expectedResult := "000001"
		data := ValuesEncoding[e]

		message := Message{}
		fieldsPackager := packager.FieldsPackager{}
		fieldsPackager.Encoding = encoding

		fields := make(map[string]packager.FieldsPackager)
		fields["011"] = fieldsPackager

		pkg := packager.Packager{}
		pkg.Fields = fields
		message.Packager = &pkg

		result, _, err := message.UnpackEncoding(data, "011", 0, 6)
		if err != nil {
			t.Fatalf(`UnpackEncoding(%s) Encoding=%s - Error %s`, data, encoding, err.Error())
		}

		if result != expectedResult {
			t.Fatalf(`UnpackEncoding(%s) Encoding=%s - Result "%s" does not match "%s"`, data, encoding, result, expectedResult)
		}
		t.Logf(`UnpackEncoding=%-28s Encoding=%-6s - Result "%-6s" match "%-6s"`, data, encoding, result, expectedResult)
	}
}

// TestPackEncoding calls encoding.PackEncoding
func TestPackEncoding(t *testing.T) {
	for e, encoding := range FieldEncoding {
		data := "000001"
		expectedResult := ValuesEncoding[e]

		message := Message{}
		fieldsPackager := packager.FieldsPackager{}
		fieldsPackager.Encoding = encoding

		fields := make(map[string]packager.FieldsPackager)
		fields["011"] = fieldsPackager

		pkg := packager.Packager{}
		pkg.Fields = fields
		message.Packager = &pkg

		message.SetField("011", data)

		result := message.PackEncoding("011", "", "")

		if result != expectedResult {
			t.Fatalf(`PackEncoding(%s) Encoding=%s - Result "%s" does not match "%s"`, data, encoding, result, expectedResult)
		}
		t.Logf(`PackEncoding=%-11s Encoding=%-6s - Result "%-6s" match "%-6s"`, data, encoding, result, expectedResult)
	}
}
