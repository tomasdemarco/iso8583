package encoding

import (
	"testing"
)

// TestEbcdicDecode calls encoding.EbcdicDecode
func TestEbcdicDecode(t *testing.T) {
	data := "f2f0f2f4"
	expectedResult := "2024"

	result, err := EbcdicDecode(data)
	if err != nil {
		t.Fatalf(`EbcdicDecode(%s) - Error %s`, data, err.Error())
	}

	if result != expectedResult {
		t.Fatalf(`EbcdicDecode(%s) - Result "%s" does not match "%s"`, data, result, expectedResult)
	}
}

// TestEbcdicEncode calls encoding.EbcdicEncode
func TestEbcdicEncode(t *testing.T) {
	data := "2024"
	expectedResult := "f2f0f2f4"

	result, err := EbcdicEncode(data)
	if err != nil {
		t.Fatalf(`EbcdicEncode(%s) - Error %s`, data, err.Error())
	}

	if result != expectedResult {
		t.Fatalf(`EbcdicEncode(%s) - Result "%s" does not match "%s"`, data, result, expectedResult)
	}
}
