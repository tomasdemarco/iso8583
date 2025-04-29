package encoding

import (
	"bytes"
	"testing"
)

// TestEbcdicDecode calls encoding.EbcdicDecode
func TestBcdDecode(t *testing.T) {
	data := []byte{0x20, 0x24}
	expectedResult := "2024"

	result, err := Decode(data)
	if err != nil {
		t.Fatalf(`Decode(%x) - Error %s`, data, err.Error())
	}

	if result != expectedResult {
		t.Fatalf(`Decode(%x) - Result "%x" does not match "%s"`, data, result, expectedResult)
	}

	t.Logf(`Decode(%x) - Result "%s" match "%s"`, data, result, expectedResult)
}

// TestEbcdicEncode calls encoding.EbcdicEncode
func TestBcdEncode(t *testing.T) {
	data := "202420"
	expectedResult := []byte{0x20, 0x24, 0x20}

	result, err := Encode(data)
	if err != nil {
		t.Fatalf(`Decode(%s) - Error %s`, data, err.Error())
	}

	if !bytes.Equal(result, expectedResult) {
		t.Fatalf(`Encode(%s) - Result "%x" does not match "%x"`, data, result, expectedResult)
	}

	t.Logf(`Encode(%x) - Result "%x" match "%x"`, data, result, expectedResult)
}
