package encoding

import (
	"bytes"
	"testing"
)

// TestBitmapDecode calls encoding.BitmapDecode
func TestBitmapDecode(t *testing.T) {
	data := []byte{0x00, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	expectedResult := []string{"011"}

	result, err := BitmapDecode(data, 1)
	if err != nil {
		t.Fatalf(`BitmapDecode(%s) - Error %s`, data, err.Error())
	}

	if len(result) != len(expectedResult) {
		t.Fatalf(`UnpackBitmap(%s) - Length is different - Result "%s" / Expected "%s"`, data, result, expectedResult)
	}

	for i, v := range result {
		if v != expectedResult[i] {
			t.Fatalf(`BitmapDecode(%s) - Result "%s" does not match "%s"`, data, result, expectedResult)

		}
	}
}

// TestBitmapEncode calls encoding.BitmapEncode
func TestBitmapEncode(t *testing.T) {
	data := []string{"001", "002", "004", "065", "126"}
	expectedResult := []byte{0xd0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04}

	result, err := BitmapEncode(data)
	if err != nil {
		t.Fatalf(`BitmapEncode(%s) - Error %s`, data, err.Error())
	}

	if !bytes.Equal(result, expectedResult) {
		t.Fatalf(`BitmapEncode(%s) - Result "%s" does not match "%s"`, data, result, expectedResult)
	}
}
