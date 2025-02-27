package emv

import (
	"reflect"
	"testing"
)

var (
	tags                 = []string{"82", "84", "95", "9A", "9C", "5F2A", "9F02", "9F03", "9F10", "9F1A", "9F26", "9F27", "9F33", "9F34", "9F36", "9F37", "9F1E"}
	tagsValues           = []string{"0000", "A0000000031010", "0000000000", "220504", "00", "0032", "000000001000", "000000000000", "06010A03A00000", "0032", "B155B92300C63C26", "80", "E0F8C8", "1F0002", "0012", "864D0C5F", "3530303130353331"}
	tagsWithFilter       = []string{"82", "84", "95", "9A", "9C", "5F2A", "9F02", "9F03", "9F10", "9F1A", "9F26", "9F27", "9F33", "9F34", "9F36", "9F37"}
	tagsValuesWithFilter = []string{"0000", "A0000000031010", "0000000000", "220504", "00", "0032", "000000001000", "000000000000", "06010A03A00000", "0032", "B155B92300C63C26", "80", "E0F8C8", "1F0002", "0012", "864D0C5F"}
	value                = "9F2608B155B92300C63C26820200009F360200129F100706010A03A000009F3303E0F8C8950500000000009F3704864D0C5F9A032205049C01009F02060000000010009F03060000000000009F2701809F34031F00025F2A0200329F1A0200328407A00000000310109F1E083530303130353331"
)

// TestUnpackEmv calls emv.Unpack
func TestUnpackEmv(t *testing.T) {
	resultExpected := make(map[string]string)
	for i, tag := range tags {
		resultExpected[tag] = tagsValues[i]
	}

	result, err := Unpack(value)
	if err != nil {
		t.Fatalf(`Unpack(%s) - Error %s`, value, err.Error())
	}

	if len(result) != len(resultExpected) {
		t.Fatalf(`Unpack(%s) - Length tags is different - Result "%s" / Expected "%s"`, value, result, resultExpected)
	}

	for key, value1 := range result {
		value2, ok := resultExpected[key]
		if !ok || !reflect.DeepEqual(value1, value2) {
			t.Fatalf(`Unpack(%s) - Result tags "%s" does not match "%s"`, value, result, resultExpected)
		}
	}

	t.Logf(`Unpack(%s) - Result "%s" match "%s"`, value, result, resultExpected)
}

// TestUnpackEmv calls emv.Unpack
func TestUnpackEmvWithFilter(t *testing.T) {
	resultExpected := make(map[string]string)
	for i, tag := range tagsWithFilter {
		resultExpected[tag] = tagsValuesWithFilter[i]
	}

	result, err := Unpack(value, tagsWithFilter...)
	if err != nil {
		t.Fatalf(`Unpack(%s) - Error %s`, value, err.Error())
	}

	if len(result) != len(resultExpected) {
		t.Fatalf(`Unpack(%s) - Length tags is different - Result "%s" / Expected "%s"`, value, result, resultExpected)
	}

	for key, value1 := range result {
		value2, ok := resultExpected[key]
		if !ok || !reflect.DeepEqual(value1, value2) {
			t.Fatalf(`Unpack(%s) - Result tags "%s" does not match "%s"`, value, result, resultExpected)
		}
	}

	t.Logf(`Unpack(%s) - Result "%s" match "%s"`, value, result, resultExpected)
}

// TestPackEmv calls emv.Pack
func TestPackEmv(t *testing.T) {
	resultExpected := "5F2A040032820400008414A0000000031010951000000000009A062205049C02009F02120000000010009F03120000000000009F101406010A03A000009F1A0400329F1E1635303031303533319F2616B155B92300C63C269F2702809F3306E0F8C89F34061F00029F360400129F3708864D0C5F"

	data := make(map[string]string)
	for i, tag := range tags {
		data[tag] = tagsValues[i]
	}

	result := Pack(data)

	if result != resultExpected {
		t.Fatalf(`Pack(%s) - Result "%s" does not match "%s"`, data, result, resultExpected)
	}
	t.Logf(`Pack(%s) - Result "%s" match "%s"`, data, result, resultExpected)
}
