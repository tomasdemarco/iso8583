package encoding

import (
	"bytes"
	"testing"
)

// TestDecodeAscii calls encoding.Decode
func TestDecodeAscii(t *testing.T) {
	data := []byte{0x32, 0x30, 0x32, 0x34}
	expectedResult := "2024"

	result, err := Decode(data)
	if err != nil {
		t.Fatalf(`Decode(%s) - Error %s`, data, err.Error())
	}

	if result != expectedResult {
		t.Fatalf(`Decode(%s) - Result "%s" does not match "%s"`, data, result, expectedResult)
	}

	t.Logf(`Decode(%s) - Result "%x" match "%x"`, data, result, expectedResult)
}

// TestEncodeAscii calls encoding.Encode
func TestEncodeAscii(t *testing.T) {
	data := "2024"
	expectedResult := []byte{0x32, 0x30, 0x32, 0x34}

	result := Encode(data)
	if !bytes.Equal(result, expectedResult) {
		t.Fatalf(`Encode(%x) - Result "%s" does not match "%s"`, data, result, expectedResult)
	}

	t.Logf(`Encode(%x) - Result "%s" match "%s"`, data, result, expectedResult)
}
