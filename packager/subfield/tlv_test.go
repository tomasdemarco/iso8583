package subfield

import (
	"reflect"
	"testing"
)

func TestTLVPacker(t *testing.T) {
	packer := &TLVPackager{}

	t.Run("Pack Success", func(t *testing.T) {
		data := map[string]string{
			"9f02": "000000001500",
			"9f03": "000000000000",
		}
		// Expected: tag(9f02) + length(06) + value(000000001500) + tag(9f03) + length(06) + value(000000000000)
		expected := "9f02060000000015009f0306000000000000"
		result, err := packer.Pack(data)
		if err != nil {
			t.Fatalf("Pack() error = %v, wantErr %v", err, false)
		}
		if result != expected {
			t.Errorf("Pack() = %v, want %v", result, expected)
		}
	})

	t.Run("Pack Error Invalid Tag", func(t *testing.T) {
		data := map[string]string{"INVALID": "value"}
		_, err := packer.Pack(data)
		if err == nil {
			t.Errorf("Pack() expected an error for invalid tag, but got nil")
		}
	})

	t.Run("Pack Error Invalid value", func(t *testing.T) {
		data := map[string]string{"9f02": "INVALID"}
		_, err := packer.Pack(data)
		if err == nil {
			t.Errorf("Pack() expected an error for invalid value, but got nil")
		}
	})
}

func TestTLVUnpacker(t *testing.T) {
	packer := &TLVPackager{}

	t.Run("Unpack Success", func(t *testing.T) {
		// Hex representation of: tag(9f02) + length(06) + value(000000001500) + tag(9f03) + length(06) + value(000000000000)
		tlvString := "9f02060000000015009f0306000000000000"
		expected := map[string]string{
			"9f02": "000000001500",
			"9f03": "000000000000",
		}
		result, err := packer.Unpack(tlvString)
		if err != nil {
			t.Fatalf("Unpack() error = %v, wantErr %v", err, false)
		}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Unpack() = %v, want %v", result, expected)
		}
	})

	t.Run("Unpack Error Malformed String", func(t *testing.T) {
		tlvString := "9f020" // Incomplete
		_, err := packer.Unpack(tlvString)
		if err == nil {
			t.Errorf("Unpack() expected an error for malformed string, but got nil")
		}
	})

	t.Run("Unpack Error Length Exceeds Buffer", func(t *testing.T) {
		tlvString := "9f020a" // Length 10, but no value
		_, err := packer.Unpack(tlvString)
		if err == nil {
			t.Errorf("Unpack() expected an error for length exceeding buffer, but got nil")
		}
	})
}

func TestTLVRoundTrip(t *testing.T) {
	packer := &TLVPackager{}
	originalData := map[string]string{
		"9a":   "251231",
		"9f02": "000000001500",
		"9f27": "80",
	}

	packed, err := packer.Pack(originalData)
	if err != nil {
		t.Fatalf("Roundtrip Pack failed: %v", err)
	}

	unpacked, err := packer.Unpack(packed)
	if err != nil {
		t.Fatalf("Roundtrip Unpack failed: %v", err)
	}

	if !reflect.DeepEqual(originalData, unpacked) {
		t.Errorf("Roundtrip data mismatch. Original: %v, Unpacked: %v", originalData, unpacked)
	}
}
