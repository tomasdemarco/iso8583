package message

import (
	"gitlab.com/g6604/adquirencia/desarrollo/golang_package/iso8583/encoding"
	"gitlab.com/g6604/adquirencia/desarrollo/golang_package/iso8583/packager"
	"testing"
)

// TestUnpackBitmap calls message.UnpackBitmap
func TestUnpackBitmap(t *testing.T) {
	expectedResult := []string{"001", "002", "004", "065", "126"}

	// Encoding BCD
	data := "d0000000000000008000000000000004"
	message := Message{}
	fieldsPackager := packager.Field{}
	fieldsPackager.Encoding = encoding.Bcd

	fields := make(map[string]packager.Field)
	fields["001"] = fieldsPackager

	pkg := packager.Packager{}
	pkg.Fields = fields
	message.Packager = &pkg

	_, result, err := message.UnpackBitmap(0, data)
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

	message.Packager.Fields["001"] = fieldsPackager

	_, result, err = message.UnpackBitmap(0, data)
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

	message.Packager.Fields["001"] = fieldsPackager

	_, result, err = message.UnpackBitmap(0, data)
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
	data := []string{"004", "011"}
	expectedResult := "1020000000000000"

	message := Message{}
	message.SetField("004", "000001000000")
	message.SetField("011", "000001")

	result, err := message.PackBitmap()
	if err != nil {
		t.Fatalf(`PackBitmap(%s) - Error %s`, data, err.Error())
	}

	if result != expectedResult {
		t.Fatalf(`PackBitmap(%s) - Result "%s" does not match "%s"`, data, result, expectedResult)
	}
}
