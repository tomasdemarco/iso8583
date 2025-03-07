package encoding

import (
	"bytes"
	"testing"
)

// TestEbcdicDecode calls encoding.EbcdicDecode
func TestEbcdicDecode(t *testing.T) {
	data := []byte{0xF2, 0xF0, 0xF2, 0xF4}
	expectedResult := "2024"

	result := EbcdicDecode(data)

	if result != expectedResult {
		t.Fatalf(`EbcdicDecode(%x) - Result "%s" does not match "%s"`, data, result, expectedResult)
	}

	t.Logf(`EbcdicDecode(%x) - Result "%s" match "%s"`, data, result, expectedResult)
}

// TestEbcdicEncode calls encoding.EbcdicEncode
func TestEbcdicEncode(t *testing.T) {
	data := "2024"
	expectedResult := []byte{0xF2, 0xF0, 0xF2, 0xF4}

	result := EbcdicEncode(data)

	if !bytes.Equal(result, expectedResult) {
		t.Fatalf(`EbcdicEncode(%s) - Result "%s" does not match "%s"`, data, result, expectedResult)
	}

	t.Logf(`EbcdicEncode(%x) - Result "%s" match "%s"`, data, result, expectedResult)
}
