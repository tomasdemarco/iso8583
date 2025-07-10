package encoding

import (
	"bytes"
	"testing"
)

func TestEncodersRoundTrip(t *testing.T) {
	testCases := []struct {
		name           string
		encoder        Encoder
		inputString    string
		expectedBytes  []byte
		expectedString string
		setLength      int // Longitud a establecer en el encoder antes de Encode/Decode
	}{
		// --- ASCII Encoder ---
		{
			name:           "ASCII - 'Hello'",
			encoder:        &ASCII{},
			inputString:    "Hello",
			expectedBytes:  []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f},
			expectedString: "Hello",
			setLength:      5,
		},
		{
			name:           "ASCII - '123456'",
			encoder:        &ASCII{},
			inputString:    "123456",
			expectedBytes:  []byte{0x31, 0x32, 0x33, 0x34, 0x35, 0x36},
			expectedString: "123456",
			setLength:      6,
		},

		// --- BCD Encoder ---
		{
			name:           "BCD - '123456'",
			encoder:        NewBcdEncoder(false), // No padLeft para este caso
			inputString:    "123456",
			expectedBytes:  []byte{0x12, 0x34, 0x56},
			expectedString: "123456",
			setLength:      3,
		},
		{
			name:           "BCD - '012345' (padLeft)",
			encoder:        NewBcdEncoder(true), // padLeft para este caso
			inputString:    "012345",
			expectedBytes:  []byte{0x01, 0x23, 0x45},
			expectedString: "012345",
			setLength:      3,
		},
		{
			name:           "BCD - '123' (padLeft)",
			encoder:        NewBcdEncoder(true), // padLeft para este caso
			inputString:    "123",
			expectedBytes:  []byte{0x01, 0x23},
			expectedString: "0123", // BCD decodifica a longitud par
			setLength:      2,
		},

		// --- BINARY Encoder (asumiendo que maneja hex strings) ---
		{
			name:           "BINARY - 'FF'",
			encoder:        NewBinaryEncoder(),
			inputString:    "FF",
			expectedBytes:  []byte{0xFF},
			expectedString: "FF",
			setLength:      1,
		},
		{
			name:           "BINARY - '0100'",
			encoder:        NewBinaryEncoder(),
			inputString:    "0100",
			expectedBytes:  []byte{0x01, 0x00},
			expectedString: "0100",
			setLength:      2,
		},

		// --- EBCDIC Encoder (ejemplo simple) ---
		{
			name:           "EBCDIC - 'Hello'",
			encoder:        &EBCDIC{},
			inputString:    "Hello",
			expectedBytes:  []byte{0xC8, 0x85, 0x93, 0x93, 0x96},
			expectedString: "Hello",
			setLength:      5,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Establecer la longitud si es necesario
			tc.encoder.SetLength(tc.setLength)

			// --- Test Encode ---
			encodedBytes, err := tc.encoder.Encode(tc.inputString)
			if err != nil {
				t.Fatalf("Encode() falló: %v", err)
			}
			if !bytes.Equal(encodedBytes, tc.expectedBytes) {
				t.Errorf("Encode() bytes incorrectos.\nEsperado: %x\nRecibido:  %x", tc.expectedBytes, encodedBytes)
			}

			// --- Test Decode ---
			decodedString, err := tc.encoder.Decode(tc.expectedBytes)
			if err != nil {
				t.Fatalf("Decode() falló: %v", err)
			}
			if decodedString != tc.expectedString {
				t.Errorf("Decode() string incorrecto.\nEsperado: %s\nRecibido:  %s", tc.expectedString, decodedString)
			}
		})
	}
}
