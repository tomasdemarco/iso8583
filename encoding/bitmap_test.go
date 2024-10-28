package encoding

import "testing"

// TestBitmapDecode calls encoding.BitmapDecode
func TestBitmapDecode(t *testing.T) {
	data := "0020000000000000"
	expectedResult := []string{"011"}

	result, err := BitmapDecode(data)
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
	expectedResult := "d0000000000000008000000000000004"

	result, err := BitmapEncode(data)
	if err != nil {
		t.Fatalf(`BitmapEncode(%s) - Error %s`, data, err.Error())
	}

	if result != expectedResult {
		t.Fatalf(`BitmapEncode(%s) - Result "%s" does not match "%s"`, data, result, expectedResult)
	}
}
