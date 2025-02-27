package bitmap

import (
	"github.com/tomasdemarco/iso8583/encoding"
	"github.com/tomasdemarco/iso8583/field"
	"github.com/tomasdemarco/iso8583/packager"
	"testing"
)

// TestUnpackBitmap calls message.UnpackBitmap
func TestUnpackBitmap(t *testing.T) {
	expectedResult := []string{"001", "002", "004", "065", "126"}

	// Encoding BCD
	data := "d0000000000000008000000000000004"

	fieldsPackager := packager.Field{}
	fieldsPackager.Encoding = encoding.Bcd

	_, result, err := Unpack(fieldsPackager, 0, data)
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

	// Encoding ASCII
	data = "6430303030303030303030303030303038303030303030303030303030303034"
	fieldsPackager = packager.Field{}
	fieldsPackager.Encoding = encoding.Ascii

	_, result, err = Unpack(fieldsPackager, 0, data)
	if err != nil {
		t.Fatalf(`UnpackBitmap(%s) - Error %s`, data, err.Error())
	}

	if len(result) != len(expectedResult) {
		t.Fatalf(`UnpackBitmap(%s) - Length is different`, data)
	}

	for i, v := range result {
		if v != expectedResult[i] {
			t.Fatalf(`UnpackBitmap(%s) - Result "%s" does not match "%s"`, data, result, expectedResult)

		}
	}

	// Encoding EBCDIC
	data = "84f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f8f0f0f0f0f0f0f0f0f0f0f0f0f0f0f4"
	fieldsPackager = packager.Field{}
	fieldsPackager.Encoding = encoding.Ebcdic

	_, result, err = Unpack(fieldsPackager, 0, data)
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
}

// TestPackBitmap calls message.PackBitmap
func TestPackBitmap(t *testing.T) {
	expectedResultArr := []string{"004", "011"}
	expectedResult := "1020000000000000"

	fields := make(map[string]field.Field)

	fieldAux := field.Field{Value: "000001000000"}
	fields["004"] = fieldAux

	fieldAux = field.Field{Value: "000001"}
	fields["011"] = fieldAux

	resultArr, resultStr, err := Pack(fields)
	if err != nil {
		t.Fatalf(`PackBitmap(%s) - Error %s`, expectedResultArr, err.Error())
	}

	if *resultStr != expectedResult {
		t.Fatalf(`PackBitmap(%s) - Result "%s" does not match "%s"`, expectedResultArr, *resultStr, expectedResult)
	}

	t.Logf(`PackBitmap(%s) - Result "%s" match "%s"`, expectedResultArr, *resultStr, expectedResult)

	if len(resultArr) != len(expectedResultArr) {
		t.Fatalf(`UnpackBitmap(%s) - Length is different - Result "%s" / Expected "%s"`, expectedResultArr, resultArr, expectedResultArr)
	}

	t.Logf(`PackBitmap(%s) - Result "%s" match "%s"`, expectedResultArr, resultArr, expectedResultArr)

	for i, v := range resultArr {
		if v != expectedResultArr[i] {
			t.Fatalf(`UnpackBitmap(%s) - Result "%s" does not match "%s"`, expectedResultArr, resultArr, expectedResultArr)
		}
	}
}
