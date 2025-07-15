package encoding

import (
	"bytes"
	"errors"
	"testing"
)

func TestEncodersRoundTrip(t *testing.T) {
	testCases := []struct {
		name            string
		encoder         Encoder
		inputString     string
		expectedBytes   []byte
		expectedString  string
		setLength       int // Longitud a establecer en el encoder antes de Encode/Decode
		expectEncodeErr bool
		expectDecodeErr bool
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
		{
			name:            "ASCII - Decode error (not enough data)",
			encoder:         &ASCII{},
			inputString:     "",
			expectedBytes:   []byte{},
			expectedString:  "",
			setLength:       5, // Espera 5 bytes, pero solo se proporcionan 2
			expectEncodeErr: false,
			expectDecodeErr: true,
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
		{
			name:           "BCD - '12D4' (decode D to =)",
			encoder:        NewBcdEncoder(false),
			inputString:    "12=4",
			expectedBytes:  []byte{0x12, 0xD4},
			expectedString: "12=4",
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
			if tc.expectEncodeErr {
				if err == nil {
					t.Fatalf("Encode() esperaba un error, pero no lo obtuvo")
				}
			} else {
				if err != nil {
					t.Fatalf("Encode() falló: %v", err)
				}
				if !bytes.Equal(encodedBytes, tc.expectedBytes) {
					t.Errorf("Encode() bytes incorrectos.\nEsperado: %x\nRecibido:  %x", tc.expectedBytes, encodedBytes)
				}
			}

			// --- Test Decode ---
			decodedString, err := tc.encoder.Decode(tc.expectedBytes)
			if tc.expectDecodeErr {
				if err == nil {
					t.Fatalf("Decode() esperaba un error, pero no lo obtuvo")
				}
				return // Prueba de error finalizada
			}

			if err != nil {
				t.Fatalf("Decode() falló: %v", err)
			}
			if decodedString != tc.expectedString {
				t.Errorf("Decode() string incorrecto.\nEsperado: %s\nRecibido:  %s", tc.expectedString, decodedString)
			}
		})
	}
}

func TestEncodingTypeMethods(t *testing.T) {
	// Test String() and EnumIndex()
	t.Run("String and EnumIndex", func(t *testing.T) {
		var typ Encoding

		typ = None
		if typ.String() != "None" || typ.EnumIndex() != 0 {
			t.Errorf("None: Expected String 'None', EnumIndex 0; Got '%s', %d", typ.String(), typ.EnumIndex())
		}

		typ = Bcd
		if typ.String() != "BCD" || typ.EnumIndex() != 1 {
			t.Errorf("Bcd: Expected String 'BCD', EnumIndex 1; Got '%s', %d", typ.String(), typ.EnumIndex())
		}

		typ = Ascii
		if typ.String() != "ASCII" || typ.EnumIndex() != 2 {
			t.Errorf("Ascii: Expected String 'ASCII', EnumIndex 2; Got '%s', %d", typ.String(), typ.EnumIndex())
		}

		typ = Ebcdic
		if typ.String() != "EBCDIC" || typ.EnumIndex() != 3 {
			t.Errorf("Ebcdic: Expected String 'EBCDIC', EnumIndex 3; Got '%s', %d", typ.String(), typ.EnumIndex())
		}

		typ = Hex
		if typ.String() != "HEX" || typ.EnumIndex() != 4 {
			t.Errorf("Hex: Expected String 'HEX', EnumIndex 4; Got '%s', %d", typ.String(), typ.EnumIndex())
		}

		typ = Binary
		if typ.String() != "BINARY" || typ.EnumIndex() != 5 {
			t.Errorf("Binary: Expected String 'BINARY', EnumIndex 5; Got '%s', %d", typ.String(), typ.EnumIndex())
		}
	})

	// Test UnmarshalJSON()
	t.Run("UnmarshalJSON", func(t *testing.T) {
		var typ Encoding

		// Success cases
		err := typ.UnmarshalJSON([]byte(`"None"`))
		if err != nil {
			t.Fatalf("UnmarshalJSON for 'None' failed: %v", err)
		}
		if typ != None {
			t.Errorf("UnmarshalJSON for 'None': Expected None, Got %v", typ)
		}

		err = typ.UnmarshalJSON([]byte(`"BCD"`))
		if err != nil {
			t.Fatalf("UnmarshalJSON for 'BCD' failed: %v", err)
		}
		if typ != Bcd {
			t.Errorf("UnmarshalJSON for 'BCD': Expected Bcd, Got %v", typ)
		}

		err = typ.UnmarshalJSON([]byte(`"ASCII"`))
		if err != nil {
			t.Fatalf("UnmarshalJSON for 'ASCII' failed: %v", err)
		}
		if typ != Ascii {
			t.Errorf("UnmarshalJSON for 'ASCII': Expected Ascii, Got %v", typ)
		}

		err = typ.UnmarshalJSON([]byte(`"EBCDIC"`))
		if err != nil {
			t.Fatalf("UnmarshalJSON for 'EBCDIC' failed: %v", err)
		}
		if typ != Ebcdic {
			t.Errorf("UnmarshalJSON for 'EBCDIC': Expected Ebcdic, Got %v", typ)
		}

		err = typ.UnmarshalJSON([]byte(`"HEX"`))
		if err != nil {
			t.Fatalf("UnmarshalJSON for 'HEX' failed: %v", err)
		}
		if typ != Hex {
			t.Errorf("UnmarshalJSON for 'HEX': Expected Hex, Got %v", typ)
		}

		err = typ.UnmarshalJSON([]byte(`"BINARY"`))
		if err != nil {
			t.Fatalf("UnmarshalJSON for 'BINARY' failed: %v", err)
		}
		if typ != Binary {
			t.Errorf("UnmarshalJSON for 'BINARY': Expected Binary, Got %v", typ)
		}

		// Error case: invalid string
		err = typ.UnmarshalJSON([]byte(`"INVALID"`))
		if err == nil {
			t.Fatalf("UnmarshalJSON for 'INVALID' expected an error, got nil")
		}
		if !errors.Is(err, ErrInvalidEncodingType) {
			t.Errorf("UnmarshalJSON for 'INVALID': Expected ErrInvalidEncodingType, Got %v", err)
		}

		// Error case: non-string input
		err = typ.UnmarshalJSON([]byte(`123`))
		if err == nil {
			t.Fatalf("UnmarshalJSON for non-string expected an error, got nil")
		}
	})

	// Test IsValid()
	t.Run("IsValid", func(t *testing.T) {
		var typ Encoding

		typ = Ascii
		if !typ.IsValid() {
			t.Errorf("Ascii: Expected IsValid true, Got false")
		}

		typ = Encoding(99) // Invalid type
		if typ.IsValid() {
			t.Errorf("Invalid Encoding: Expected IsValid false, Got true")
		}
	})
}
