package prefix

import (
	"bytes"
	"errors"
	"testing"
)

func TestPrefixersRoundTrip(t *testing.T) {
	ll := LL
	lll := LLL
	testCases := []struct {
		name            string
		prefixer        Prefixer
		length          int
		expectedBytes   []byte
		expectEncodeErr bool
	}{
		// --- ASCII Prefixer Cases ---
		{"ASCII LL, len 7", NewAsciiPrefixer(ll.EnumIndex(), false, false), 7, []byte{'0', '7'}, false},
		{"ASCII LL, len 99", NewAsciiPrefixer(ll.EnumIndex(), false, false), 99, []byte{'9', '9'}, false},
		{"ASCII LL, len 100 (error)", NewAsciiPrefixer(ll.EnumIndex(), false, false), 100, nil, true},
		{"ASCII LLL, len 123", NewAsciiPrefixer(lll.EnumIndex(), false, false), 123, []byte{'1', '2', '3'}, false},
		{"ASCII LLL, len 999", NewAsciiPrefixer(lll.EnumIndex(), false, false), 999, []byte{'9', '9', '9'}, false},
		{"ASCII LLL, len 1000 (error)", NewAsciiPrefixer(lll.EnumIndex(), false, false), 1000, nil, true},
		{"ASCII LL, hex true, inclusive true, len 10", NewAsciiPrefixer(ll.EnumIndex(), true, true), 10, []byte{0x30, 0x63}, false},

		// --- BCD Prefixer Cases ---
		{"BCD LL, len 7", NewBcdPrefixer(ll.EnumIndex(), false, false), 7, []byte{0x07}, false},
		{"BCD LL, len 12", NewBcdPrefixer(ll.EnumIndex(), false, false), 12, []byte{0x12}, false},
		{"BCD LL, len 99", NewBcdPrefixer(ll.EnumIndex(), false, false), 99, []byte{0x99}, false},
		{"BCD LLL, len 123", NewBcdPrefixer(lll.EnumIndex(), false, false), 123, []byte{0x01, 0x23}, false},
		{"BCD LLL, len 999", NewBcdPrefixer(lll.EnumIndex(), false, false), 999, []byte{0x09, 0x99}, false},
		{"BCD LL, len 100 (error)", NewBcdPrefixer(ll.EnumIndex(), false, false), 100, nil, true},
		{"BCD LL, hex true, inclusive true, len 10", NewBcdPrefixer(ll.EnumIndex(), false, true), 10, []byte{0x12}, false},

		// --- Binary Prefixer Cases ---
		{"Binary LL, len 255", NewBinaryPrefixer(ll.EnumIndex(), false), 255, []byte{0xff}, false},
		{"Binary LL, len 256 (error)", NewBinaryPrefixer(ll.EnumIndex(), false), 256, nil, true},
		{"Binary LLL, len 256", NewBinaryPrefixer(lll.EnumIndex(), false), 256, []byte{0x01, 0x00}, false},
		{"Binary LLL, len 65535", NewBinaryPrefixer(lll.EnumIndex(), false), 65535, []byte{0xff, 0xff}, false},
		{"Binary LLL, len 65536 (error)", NewBinaryPrefixer(lll.EnumIndex(), false), 65536, nil, true},
		{"Binary LL, inclusive true, len 10", NewBinaryPrefixer(ll.EnumIndex(), true), 10, []byte{0x0B}, false},

		// --- EBCDIC Prefixer Cases ---
		{"EBCDIC LL, len 7", NewEbcdicPrefixer(ll.EnumIndex(), false, false), 7, []byte{0xF0, 0xF7}, false},
		{"EBCDIC LL, len 99", NewEbcdicPrefixer(ll.EnumIndex(), false, false), 99, []byte{0xF9, 0xF9}, false},
		{"EBCDIC LL, len 100 (error)", NewEbcdicPrefixer(ll.EnumIndex(), false, false), 100, nil, true},
		{"EBCDIC LLL, len 123", NewEbcdicPrefixer(lll.EnumIndex(), false, false), 123, []byte{0xF1, 0xF2, 0xF3}, false},
		{"EBCDIC LLL, len 999", NewEbcdicPrefixer(lll.EnumIndex(), false, false), 999, []byte{0xF9, 0xF9, 0xF9}, false},
		{"EBCDIC LLL, len 1000 (error)", NewEbcdicPrefixer(lll.EnumIndex(), false, false), 1000, nil, true},
		{"EBCDIC LL, hex true, inclusive true, len 10", NewEbcdicPrefixer(ll.EnumIndex(), false, true), 10, []byte{0xF1, 0xF2}, false},

		// --- None Prefixer Cases ---
		{"None Fixed, len 0", NewNonePrefixer(0), 0, nil, false},
		{"None Fixed, len 10", NewNonePrefixer(10), 0, nil, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// --- Test EncodeLength ---
			encoded, err := tc.prefixer.EncodeLength(tc.length)

			if tc.expectEncodeErr {
				if err == nil {
					t.Fatalf("EncodeLength() esperaba un error, pero no lo obtuvo")
				}
				return // Prueba de error finalizada
			}

			if err != nil {
				t.Fatalf("EncodeLength() falló: %v", err)
			}

			if !bytes.Equal(tc.expectedBytes, encoded) {
				t.Fatalf("EncodeLength() bytes incorrectos. Esperado: %x, Recibido: %x", tc.expectedBytes, encoded)
			}

			// --- Test DecodeLength ---
			decoded, err := tc.prefixer.DecodeLength(encoded, 0)
			if err != nil {
				t.Fatalf("DecodeLength() falló: %v", err)
			}

			if tc.length != decoded {
				t.Fatalf("DecodeLength() longitud incorrecta. Esperado: %d, Recibido: %d", tc.length, decoded)
			}
		})
	}
}

func TestPrefixerDecodeErrors(t *testing.T) {
	ll := LL
	lll := LLL

	testCases := []struct {
		name        string
		prefixer    Prefixer
		data        []byte
		expectedErr error
	}{
		{
			name:        "ASCII LL - Not enough data for decode",
			prefixer:    NewAsciiPrefixer(ll.EnumIndex(), false, false),
			data:        []byte{'0'},
			expectedErr: ErrFailedToDecodeLength,
		},
		{
			name:        "ASCII LLL - Invalid length string conversion (non-numeric)",
			prefixer:    NewAsciiPrefixer(lll.EnumIndex(), false, false),
			data:        []byte{'A', 'B', 'C'},
			expectedErr: ErrInvalidLengthStringConversion,
		},
		{
			name:        "ASCII LL - Invalid hex string for lenStrToInt",
			prefixer:    NewAsciiPrefixer(ll.EnumIndex(), true, false),
			data:        []byte{'1', 'G'},
			expectedErr: ErrInvalidLengthStringConversion,
		},
		{
			name:        "ASCII LL - Invalid non-hex string for lenStrToInt",
			prefixer:    NewAsciiPrefixer(ll.EnumIndex(), false, false),
			data:        []byte{'A', 'B'},
			expectedErr: ErrInvalidLengthStringConversion,
		},
		{
			name:        "BCD LL - Not enough data for decode",
			prefixer:    NewBcdPrefixer(ll.EnumIndex(), false, false),
			data:        []byte{},
			expectedErr: ErrFailedToDecodeLength,
		},
		{
			name:        "BCD LLL - Invalid length string conversion (non-numeric)",
			prefixer:    NewBcdPrefixer(lll.EnumIndex(), false, false),
			data:        []byte{0x12, 0x3A},
			expectedErr: ErrInvalidLengthStringConversion,
		},
		{
			name:        "Binary LL - Not enough data for decode",
			prefixer:    NewBinaryPrefixer(ll.EnumIndex(), false),
			data:        []byte{},
			expectedErr: ErrFailedToDecodeLength,
		},
		{
			name:        "EBCDIC LL - Not enough data for decode",
			prefixer:    NewEbcdicPrefixer(ll.EnumIndex(), false, false),
			data:        []byte{0xF0},
			expectedErr: ErrFailedToDecodeLength,
		},
		{
			name:        "EBCDIC LLL - Invalid length string conversion (non-numeric)",
			prefixer:    NewEbcdicPrefixer(lll.EnumIndex(), false, false),
			data:        []byte{0xF1, 0xF2, 0x40},
			expectedErr: ErrInvalidLengthStringConversion,
		},
		{
			name:        "ASCII LL Inclusive - Decoded length less than nDigits",
			prefixer:    NewAsciiPrefixer(ll.EnumIndex(), false, true),
			data:        []byte{'0', '1'},
			expectedErr: ErrInvalidLengthStringConversion,
		},
		{
			name:        "BCD LL Inclusive - Decoded length less than nDigits",
			prefixer:    NewBcdPrefixer(ll.EnumIndex(), false, true),
			data:        []byte{0x01},
			expectedErr: ErrInvalidLengthStringConversion,
		},
		{
			name:        "Binary LL Inclusive - Decoded length less than nBytes",
			prefixer:    NewBinaryPrefixer(4, true),
			data:        []byte{0x00, 0x01},
			expectedErr: ErrInvalidLengthStringConversion,
		},
		{
			name:        "EBCDIC LL Inclusive - Decoded length less than nDigits",
			prefixer:    NewEbcdicPrefixer(ll.EnumIndex(), false, true),
			data:        []byte{0xF0, 0xF1},
			expectedErr: ErrInvalidLengthStringConversion,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.prefixer.DecodeLength(tc.data, 0)

			if err == nil {
				t.Fatalf("DecodeLength() esperaba un error, pero no lo obtuvo")
			}
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("DecodeLength() error incorrecto. Esperado: %v, Recibido: %v", tc.expectedErr, err)
			}
		})
	}
}

func TestPrefixTypeMethods(t *testing.T) {
	// Test String() and EnumIndex()
	t.Run("String and EnumIndex", func(t *testing.T) {
		var typ Type

		typ = Fixed
		if typ.String() != "FIXED" || typ.EnumIndex() != 0 {
			t.Errorf("Fixed: Expected String 'FIXED', EnumIndex 0; Got '%s', %d", typ.String(), typ.EnumIndex())
		}

		typ = LL
		if typ.String() != "LL" || typ.EnumIndex() != 2 {
			t.Errorf("LL: Expected String 'LL', EnumIndex 2; Got '%s', %d", typ.String(), typ.EnumIndex())
		}
	})

	// Test UnmarshalJSON()
	t.Run("UnmarshalJSON", func(t *testing.T) {
		var typ Type

		// Success case
		err := typ.UnmarshalJSON([]byte(`"L"`))
		if err != nil {
			t.Fatalf("UnmarshalJSON for 'L' failed: %v", err)
		}
		if typ != L {
			t.Errorf("UnmarshalJSON for 'L': Expected L, Got %v", typ)
		}

		// Error case: invalid string
		err = typ.UnmarshalJSON([]byte(`"INVALID"`))
		if err == nil {
			t.Fatalf("UnmarshalJSON for 'INVALID' expected an error, got nil")
		}
		if !errors.Is(err, ErrInvalidPrefixType) {
			t.Errorf("UnmarshalJSON for 'INVALID': Expected ErrInvalidPrefixType, Got %v", err)
		}

		// Error case: non-string input
		err = typ.UnmarshalJSON([]byte(`123`))
		if err == nil {
			t.Fatalf("UnmarshalJSON for non-string expected an error, got nil")
		}
	})

	// Test IsValid()
	t.Run("IsValid", func(t *testing.T) {
		var typ Type

		typ = LL
		if !typ.IsValid() {
			t.Errorf("LL: Expected IsValid true, Got false")
		}

		typ = Type(99) // Invalid type
		if typ.IsValid() {
			t.Errorf("Invalid Type: Expected IsValid false, Got true")
		}
	})
}

func TestNewPrefixers(t *testing.T) {
	t.Run("NewAsciiPrefixer", func(t *testing.T) {
		prefixer := NewAsciiPrefixer(2, true, true).(*AsciiPrefixer)
		if prefixer.nDigits != 2 || prefixer.hex != true || prefixer.isInclusive != true {
			t.Errorf("NewAsciiPrefixer: Expected nDigits 2, hex true, isInclusive true; Got %d, %t, %t", prefixer.nDigits, prefixer.hex, prefixer.isInclusive)
		}
		if prefixer.GetPackedLength() != 2 {
			t.Errorf("NewAsciiPrefixer GetPackedLength: Expected 2, Got %d", prefixer.GetPackedLength())
		}
	})

	t.Run("NewBcdPrefixer", func(t *testing.T) {
		prefixer := NewBcdPrefixer(4, false, true).(*BcdPrefixer)
		if prefixer.nDigits != 4 || prefixer.hex != false || prefixer.isInclusive != true {
			t.Errorf("NewBcdPrefixer: Expected nDigits 4, hex false, isInclusive true; Got %d, %t, %t", prefixer.nDigits, prefixer.hex, prefixer.isInclusive)
		}
		if prefixer.GetPackedLength() != 2 {
			t.Errorf("NewBcdPrefixer GetPackedLength: Expected 2, Got %d", prefixer.GetPackedLength())
		}
	})

	t.Run("NewBinaryPrefixer", func(t *testing.T) {
		prefixer := NewBinaryPrefixer(1, true).(*BinaryPrefixer)
		if prefixer.nBytes != 1 || prefixer.isInclusive != true {
			t.Errorf("NewBinaryPrefixer: Expected nBytes 1, isInclusive true; Got %d, %t", prefixer.nBytes, prefixer.isInclusive)
		}
		if prefixer.GetPackedLength() != 1 {
			t.Errorf("NewBinaryPrefixer GetPackedLength: Expected 1, Got %d", prefixer.GetPackedLength())
		}
	})

	t.Run("NewEbcdicPrefixer", func(t *testing.T) {
		prefixer := NewEbcdicPrefixer(3, true, false).(*EbcdicPrefixer)
		if prefixer.nDigits != 3 || prefixer.hex != true || prefixer.isInclusive != false {
			t.Errorf("NewEbcdicPrefixer: Expected nDigits 3, hex true, isInclusive false; Got %d, %t, %t", prefixer.nDigits, prefixer.hex, prefixer.isInclusive)
		}
		if prefixer.GetPackedLength() != 3 {
			t.Errorf("NewEbcdicPrefixer GetPackedLength: Expected 3, Got %d", prefixer.GetPackedLength())
		}
	})

	t.Run("NewNonePrefixer", func(t *testing.T) {
		prefixer := NewNonePrefixer(5).(*NonePrefixer)
		if prefixer.nDigits != 5 {
			t.Errorf("NewNonePrefixer: Expected nDigits 5; Got %d", prefixer.nDigits)
		}
		if prefixer.GetPackedLength() != 5 {
			t.Errorf("NewNonePrefixer GetPackedLength: Expected 5, Got %d", prefixer.GetPackedLength())
		}
	})
}
