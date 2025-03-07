package encoding

import (
	"bytes"
	"testing"
)

// TestHexDecode calls encoding.HexDecode
func TestHexDecode(t *testing.T) {
	data := []byte{0x10}
	expectedResult := "16"

	result, err := HexDecode(data)
	if err != nil {
		t.Fatalf(`HexDecode(%s) - Error %s`, data, err.Error())
	}

	if result != expectedResult {
		t.Fatalf(`HexDecode(%x) - Result "%s" does not match "%s"`, data, result, expectedResult)
	}
}

// TestHexEncode calls encoding.HexEncode
func TestHexEncode(t *testing.T) {
	data := "16"
	expectedResult := []byte{0x31, 0x30}

	result, err := HexEncode(data)
	if err != nil {
		t.Fatalf(`HexEncode(%s) - Error %s`, data, err.Error())
	}

	if !bytes.Equal(result, expectedResult) {
		t.Fatalf(`HexEncode(%s) - Result "%s" does not match "%s"`, data, result, expectedResult)
	}
}
