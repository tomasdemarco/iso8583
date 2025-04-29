package encoding

import (
	"bytes"
	"testing"
)

var (
	Encodings      = []Encoding{Bcd, Ascii, Ebcdic}
	ValuesEncoding = [][]byte{{0x00, 0x00, 0x01}, {0x30, 0x30, 0x30, 0x30, 0x30, 0x31}, {0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF1}}
)

// TestUnpackEncoding calls message.UnpackEncoding
func TestUnpackEncoding(t *testing.T) {
	for e, enc := range Encodings {
		data := ValuesEncoding[e]
		expectedResult := "000001"

		result, err := Unpack(enc, data)
		if err != nil {
			t.Fatalf(`UnpackEncoding(%x) Encoding=%s - Error %s`, data, enc.String(), err.Error())
		}

		if result != expectedResult {
			t.Fatalf(`UnpackEncoding(%x) Encoding=%s - Result "%s" does not match "%s"`, data, enc.String(), result, expectedResult)
		}
		t.Logf(`UnpackEncoding(%x) Encoding=%s - Result "%s" match "%s"`, data, enc.String(), result, expectedResult)
	}
}

// TestPackEncoding calls encoding.PackEncoding
func TestPackEncoding(t *testing.T) {
	for e, enc := range Encodings {
		data := "000001"
		expectedResult := ValuesEncoding[e]

		result, err := Pack(enc, data)
		if err != nil {
			t.Fatalf(`PackEncoding(%s) - Error %s`, data, err.Error())
		}

		if !bytes.Equal(result, expectedResult) {
			t.Fatalf(`PackEncoding(%s) Encoding=%s - Result "%x" does not match "%x"`, data, enc.String(), result, expectedResult)
		}

		t.Logf(`PackEncoding(%s) Encoding=%s - Result "%x" match "%x"`, data, enc.String(), result, expectedResult)
	}
}
