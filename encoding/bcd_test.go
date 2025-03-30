package encoding

import (
	"bytes"
	"testing"
)

// TestEbcdicDecode calls encoding.EbcdicDecode
func TestBcdDecode(t *testing.T) {
	data := []byte{0x20, 0x24}
	expectedResult := "2024"

	result, err := BcdDecode(data)
	if err != nil {
		t.Fatalf(`BcdDecode(%x) - Error %s`, data, err.Error())
	}

	if result != expectedResult {
		t.Fatalf(`BcdDecode(%x) - Result "%x" does not match "%s"`, data, result, expectedResult)
	}

	t.Logf(`BcdDecode(%x) - Result "%s" match "%s"`, data, result, expectedResult)
}

// TestEbcdicEncode calls encoding.EbcdicEncode
func TestBcdEncode(t *testing.T) {
	data := "202420"
	expectedResult := []byte{0x20, 0x24, 0x20}

	result, err := BcdEncode(data)
	if err != nil {
		t.Fatalf(`BcdDecode(%s) - Error %s`, data, err.Error())
	}

	if !bytes.Equal(result, expectedResult) {
		t.Fatalf(`BcdEncode(%s) - Result "%x" does not match "%x"`, data, result, expectedResult)
	}

	t.Logf(`BcdEncode(%x) - Result "%x" match "%x"`, data, result, expectedResult)
}
