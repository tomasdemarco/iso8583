package prefix

import (
	"bytes"
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

		// --- BCD Prefixer Cases ---
		{"BCD LL, len 7", NewBcdPrefixer(ll.EnumIndex(), false, false), 7, []byte{0x07}, false},
		{"BCD LL, len 12", NewBcdPrefixer(ll.EnumIndex(), false, false), 12, []byte{0x12}, false},
		{"BCD LL, len 99", NewBcdPrefixer(ll.EnumIndex(), false, false), 99, []byte{0x99}, false},
		{"BCD LLL, len 123", NewBcdPrefixer(lll.EnumIndex(), false, false), 123, []byte{0x01, 0x23}, false},
		{"BCD LLL, len 999", NewBcdPrefixer(lll.EnumIndex(), false, false), 999, []byte{0x09, 0x99}, false},

		// --- Binary Prefixer Cases ---
		{"Binary LL, len 255", NewBinaryPrefixer(ll.EnumIndex(), false), 255, []byte{0xff}, false},
		{"Binary LL, len 256 (error)", NewBinaryPrefixer(ll.EnumIndex(), false), 256, nil, true},
		{"Binary LLL, len 256", NewBinaryPrefixer(lll.EnumIndex(), false), 256, []byte{0x01, 0x00}, false},
		{"Binary LLL, len 65535", NewBinaryPrefixer(lll.EnumIndex(), false), 65535, []byte{0xff, 0xff}, false},
		{"Binary LLL, len 65536 (error)", NewBinaryPrefixer(lll.EnumIndex(), false), 65536, nil, true},
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
