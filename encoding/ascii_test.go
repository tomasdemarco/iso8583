package encoding

import (
	"testing"
)

// TestDecodeAscii calls encoding.AsciiDecode
func TestDecodeAscii(t *testing.T) {
	data := "32303234"
	expectedResult := "2024"

	result, err := AsciiDecode(data)
	if err != nil {
		t.Fatalf(`AsciiDecode(%s) - Error %s`, data, err.Error())
	}

	if result != expectedResult {
		t.Fatalf(`AsciiDecode(%s) - Result "%s" does not match "%s"`, data, result, expectedResult)
	}
}

// TestEncodeAscii calls encoding.AsciiEncode
func TestEncodeAscii(t *testing.T) {
	data := "2024"
	expectedResult := "32303234"

	result := AsciiEncode(data)
	if result != expectedResult {
		t.Fatalf(`AsciiEncode(%s) - Result "%s" does not match "%s"`, data, result, expectedResult)
	}
}
