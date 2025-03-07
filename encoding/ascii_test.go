package encoding

import (
	"bytes"
	"testing"
)

// TestDecodeAscii calls encoding.AsciiDecode
func TestDecodeAscii(t *testing.T) {
	data := []byte{0x32, 0x30, 0x32, 0x34}
	expectedResult := "2024"

	result, err := AsciiDecode(data)
	if err != nil {
		t.Fatalf(`AsciiDecode(%s) - Error %s`, data, err.Error())
	}

	if result != expectedResult {
		t.Fatalf(`AsciiDecode(%s) - Result "%s" does not match "%s"`, data, result, expectedResult)
	}

	t.Logf(`AsciiDecode(%s) - Result "%x" match "%x"`, data, result, expectedResult)
}

// TestEncodeAscii calls encoding.AsciiEncode
func TestEncodeAscii(t *testing.T) {
	data := "2024"
	expectedResult := []byte{0x32, 0x30, 0x32, 0x34}

	result := AsciiEncode(data)
	if !bytes.Equal(result, expectedResult) {
		t.Fatalf(`AsciiEncode(%x) - Result "%s" does not match "%s"`, data, result, expectedResult)
	}

	t.Logf(`AsciiEncode(%x) - Result "%s" match "%s"`, data, result, expectedResult)
}
