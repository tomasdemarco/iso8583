package bitmap

import (
	"bytes"
	"github.com/tomasdemarco/iso8583/encoding"
	pkgField "github.com/tomasdemarco/iso8583/packager/field"
	"testing"
)

// TestUnpackBitmap calls message.UnpackBitmap
func TestUnpackBitmap(t *testing.T) {
	expectedResult := []string{"001", "002", "004", "065", "126"}

	// Encoding Binary
	data := []byte{0xd0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04}

	fieldsPackager := pkgField.Field{}
	fieldsPackager.Length = 16
	fieldsPackager.Encoding = encoding.Binary

	_, result, err := Unpack(fieldsPackager, data, 0)
	if err != nil {
		t.Fatalf(`UnpackBitmap(%x) - Error %s`, data, err.Error())
	}

	if len(result) != len(expectedResult) {
		t.Fatalf(`UnpackBitmap(%x) - Length is different - Result "%s" / Expected "%s"`, data, result, expectedResult)
	}

	for i, v := range result {
		if v != expectedResult[i] {
			t.Fatalf(`UnpackBitmap(%x) - Result "%s" does not match "%s"`, data, result, expectedResult)
		}
	}

	t.Logf(`UnpackBitmap(%x) - Result "%s" match "%s"`, data, result, expectedResult)

	// Encoding ASCII
	data = []byte{0x64, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x38, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x34}

	fieldsPackager = pkgField.Field{}
	fieldsPackager.Length = 16
	fieldsPackager.Encoding = encoding.Ascii

	_, result, err = Unpack(fieldsPackager, data, 0)
	if err != nil {
		t.Fatalf(`UnpackBitmap(%x) - Error %s`, data, err.Error())
	}

	if len(result) != len(expectedResult) {
		t.Fatalf(`UnpackBitmap(%x) - Length is different`, data)
	}

	for i, v := range result {
		if v != expectedResult[i] {
			t.Fatalf(`UnpackBitmap(%x) - Result "%s" does not match "%s"`, data, result, expectedResult)

		}
	}

	t.Logf(`UnpackBitmap(%x) - Result "%s" match "%s"`, data, result, expectedResult)

	// Encoding EBCDIC
	data = []byte{0x84, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF8, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF4}

	fieldsPackager = pkgField.Field{}
	fieldsPackager.Length = 16
	fieldsPackager.Encoding = encoding.Ebcdic

	_, result, err = Unpack(fieldsPackager, data, 0)
	if err != nil {
		t.Fatalf(`UnpackBitmap(%s) - Error %s`, data, err.Error())
	}

	if len(result) != len(expectedResult) {
		t.Fatalf(`UnpackBitmap(%s) - Length is different - Result "%s" / Expected "%s"`, data, result, expectedResult)
	}

	for i, v := range result {
		if v != expectedResult[i] {
			t.Fatalf(`UnpackBitmap(%s) - Result "%s" does not match "%s"`, data, result, expectedResult)
		}
	}

	t.Logf(`UnpackBitmap(%x) - Result "%s" match "%s"`, data, result, expectedResult)
}

// TestPackBitmap calls message.PackBitmap
func TestPackBitmap(t *testing.T) {
	expectedResultArr := []string{"004", "011"}
	expectedResult := []byte{0x10, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	fields := make(map[string]string)

	fields["004"] = "000001000000"
	fields["011"] = "000001"

	resultArr, result, err := Pack(fields)
	if err != nil {
		t.Fatalf(`PackBitmap(%s) - Error %s`, expectedResultArr, err.Error())
	}

	if !bytes.Equal(result, expectedResult) {
		t.Fatalf(`PackBitmap(%s) - Result "%s" does not match "%s"`, expectedResultArr, result, expectedResult)
	}

	t.Logf(`PackBitmap(%s) - Result "%x" match "%x"`, expectedResultArr, result, expectedResult)

	if len(resultArr) != len(expectedResultArr) {
		t.Fatalf(`UnpackBitmap(%s) - Length is different - Result "%s" / Expected "%s"`, expectedResultArr, resultArr, expectedResultArr)
	}

	for i, v := range resultArr {
		if v != expectedResultArr[i] {
			t.Fatalf(`UnpackBitmap(%s) - Result "%s" does not match "%s"`, expectedResultArr, resultArr, expectedResultArr)
		}
	}

	t.Logf(`PackBitmap(%s) - Result Arr "%s" match "%s"`, expectedResultArr, resultArr, expectedResultArr)
}
