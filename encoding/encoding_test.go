package encoding

import (
	"testing"
)

var (
	Encodings      = []Encoding{Bcd, Ascii, Ebcdic}
	ValuesEncoding = []string{"000001", "303030303031", "f0f0f0f0f0f1"}
)

// TestUnpackEncoding calls message.UnpackEncoding
func TestUnpackEncoding(t *testing.T) {
	for e, enc := range Encodings {
		expectedResult := "000001"
		data := ValuesEncoding[e]

		result, _, err := Unpack(enc, data, "011", 0, 6)
		if err != nil {
			t.Fatalf(`UnpackEncoding(%s) Encoding=%s - Error %s`, data, enc.String(), err.Error())
		}

		if result != expectedResult {
			t.Fatalf(`UnpackEncoding(%s) Encoding=%s - Result "%s" does not match "%s"`, data, enc.String(), result, expectedResult)
		}
		t.Logf(`UnpackEncoding=%-28s Encoding=%-6s - Result "%-6s" match "%-6s"`, data, enc.String(), result, expectedResult)
	}
}

// TestPackEncoding calls encoding.PackEncoding
func TestPackEncoding(t *testing.T) {
	for e, enc := range Encodings {
		data := "000001"
		expectedResult := ValuesEncoding[e]

		result := Pack(enc, data)

		if result != expectedResult {
			t.Fatalf(`PackEncoding(%s) Encoding=%s - Result "%s" does not match "%s"`, data, enc.String(), result, expectedResult)
		}
		t.Logf(`PackEncoding=%-11s Encoding=%-6s - Result "%-6s" match "%-6s"`, data, enc.String(), result, expectedResult)
	}
}
