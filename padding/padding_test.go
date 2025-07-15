package padding

import (
	"testing"

	"errors"
	"github.com/tomasdemarco/iso8583/encoding"
)

func TestPaddersRoundTrip(t *testing.T) {
	// Usamos un encoder ASCII simple para la mayoría de las pruebas.
	// Para BCD, se usará un encoder BCD.
	asciiEncoder := &encoding.ASCII{}
	bcdEncoder := encoding.NewBcdEncoder(true) // Asumimos padLeft para BCD

	testCases := []struct {
		name             string
		padder           Padder
		fieldLength      int              // Longitud total esperada del campo
		dataLength       int              // Longitud de los datos sin padding
		encoder          encoding.Encoder // Encoder para simular el comportamiento de EncodePad
		expectedLeftPad  int
		expectedRightPad int
		expectError      bool
	}{
		// --- FillPadder (Left) ---
		{
			name:             "FillPadder Left - Pad 0s",
			padder:           NewFillPadder(true, "0"),
			fieldLength:      10,
			dataLength:       5,
			encoder:          asciiEncoder,
			expectedLeftPad:  5,
			expectedRightPad: 0,
		},
		{
			name:             "FillPadder Left - No Pad",
			padder:           NewFillPadder(true, "0"),
			fieldLength:      5,
			dataLength:       5,
			encoder:          asciiEncoder,
			expectedLeftPad:  0,
			expectedRightPad: 0,
		},
		{
			name:        "FillPadder Left - Data too long",
			padder:      NewFillPadder(true, "0"),
			fieldLength: 3,
			dataLength:  5,
			encoder:     asciiEncoder,
			expectError: true,
		},
		{
			name:             "FillPadder Left - Pad 0s with BCD encoder",
			padder:           NewFillPadder(true, "0"),
			fieldLength:      5, // 5 bytes BCD = 10 dígitos
			dataLength:       3, // 3 bytes BCD = 6 dígitos
			encoder:          bcdEncoder,
			expectedLeftPad:  7, // (10 - 6) / 1 = 4
			expectedRightPad: 0,
		},

		// --- FillPadder (Right) ---
		{
			name:             "FillPadder Right - Pad spaces",
			padder:           NewFillPadder(false, " "),
			fieldLength:      10,
			dataLength:       5,
			encoder:          asciiEncoder,
			expectedLeftPad:  0,
			expectedRightPad: 5,
		},

		// --- ParityPadder (Left) ---
		{
			name:             "ParityPadder Left - Odd to Even",
			padder:           NewParityPadder(true, "0"),
			fieldLength:      5, // Longitud total del campo (impar)
			dataLength:       5, // Longitud de datos (impar)
			encoder:          asciiEncoder,
			expectedLeftPad:  1,
			expectedRightPad: 0,
		},
		{
			name:             "ParityPadder Left - Even to Even (No Pad)",
			padder:           NewParityPadder(true, "0"),
			fieldLength:      6,
			dataLength:       6,
			encoder:          asciiEncoder,
			expectedLeftPad:  0,
			expectedRightPad: 0,
		},
		{
			name:             "ParityPadder Left - BCD Odd to Even",
			padder:           NewParityPadder(true, "0"),
			fieldLength:      6, // Longitud total del campo (par)
			dataLength:       5, // Longitud de datos (impar)
			encoder:          bcdEncoder,
			expectedLeftPad:  1,
			expectedRightPad: 0,
		},
		{
			name:             "ParityPadder Left - Odd to Even (Left Pad)",
			padder:           NewParityPadder(true, "X"),
			fieldLength:      6,
			dataLength:       5,
			encoder:          asciiEncoder,
			expectedLeftPad:  1,
			expectedRightPad: 0,
		},

		// --- ParityPadder (Right) ---
		{
			name:             "ParityPadder Right - Odd to Even",
			padder:           NewParityPadder(false, "0"),
			fieldLength:      6,
			dataLength:       5,
			encoder:          asciiEncoder,
			expectedLeftPad:  0,
			expectedRightPad: 1,
		},

		// --- NonePadder ---
		{
			name:             "NonePadder - No Pad",
			padder:           NONE.NONE,
			fieldLength:      5,
			dataLength:       5,
			encoder:          asciiEncoder,
			expectedLeftPad:  0,
			expectedRightPad: 0,
		},
		{
			name:        "NonePadder - Data too long (no error expected)",
			padder:      NONE.NONE,
			fieldLength: 3,
			dataLength:  5,
			encoder:     asciiEncoder,
			expectError: false, // Corregido: No se espera error aquí
		},
		{
			name:             "NonePadder - With Char",
			padder:           &NonePadder{char: "X"},
			fieldLength:      0,
			dataLength:       0,
			encoder:          asciiEncoder,
			expectedLeftPad:  0,
			expectedRightPad: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// --- Test EncodePad ---
			leftPadStr, rightPadStr, err := tc.padder.EncodePad(tc.fieldLength, tc.dataLength, tc.encoder)
			if tc.expectError {
				if err == nil {
					t.Fatalf("EncodePad() esperaba un error, pero no lo obtuvo")
				}
				return // Prueba de error finalizada
			}

			if err != nil {
				t.Fatalf("EncodePad() falló: %v", err)
			}

			// Convertir los strings de padding a int para la comparación
			leftPad := len(leftPadStr)
			rightPad := len(rightPadStr)

			if leftPad != tc.expectedLeftPad || rightPad != tc.expectedRightPad {
				t.Errorf("EncodePad() padding incorrecto. Esperado: (%d, %d), Recibido: (%d, %d)",
					tc.expectedLeftPad, tc.expectedRightPad, leftPad, rightPad)
			}

			// --- Test DecodePad ---
			// Para DecodePad, la longitud de entrada es la longitud total del campo.
			// El resultado esperado es la cantidad de padding que se eliminaría.
			decodedLeftPad, decodedRightPad := tc.padder.DecodePad(tc.fieldLength)

			// La lógica de DecodePad es más simple, solo verifica si se debe quitar padding.
			// Para FillPadder, siempre es 0,0. Para ParityPadder, es 1,0 o 0,1 si la longitud es impar.
			// Ajustamos la expectativa de DecodePad basándonos en el tipo de padder.
			var expectedDecodedLeftPad, expectedDecodedRightPad int
			switch tc.padder.(type) {
			case *FillPadder:
				expectedDecodedLeftPad = 0
				expectedDecodedRightPad = 0
			case *ParityPadder:
				if tc.fieldLength%2 != 0 { // Si la longitud del campo es impar, se espera 1 de padding.
					if tc.padder.(*ParityPadder).left {
						expectedDecodedLeftPad = 1
					} else {
						expectedDecodedRightPad = 1
					}
				}
			case *NonePadder:
				expectedDecodedLeftPad = 0
				expectedDecodedRightPad = 0
			}

			if decodedLeftPad != expectedDecodedLeftPad || decodedRightPad != expectedDecodedRightPad {
				t.Errorf("DecodePad() padding incorrecto. Esperado: (%d, %d), Recibido: (%d, %d)",
					expectedDecodedLeftPad, expectedDecodedRightPad, decodedLeftPad, decodedRightPad)
			}
		})
	}
}

func TestPaddingTypeMethods(t *testing.T) {
	// Test String() and EnumIndex()
	t.Run("String and EnumIndex", func(t *testing.T) {
		var typ Type

		typ = None
		if typ.String() != "NONE" || typ.EnumIndex() != 0 {
			t.Errorf("None: Expected String 'NONE', EnumIndex 0; Got '%s', %d", typ.String(), typ.EnumIndex())
		}

		typ = Fill
		if typ.String() != "FILL" || typ.EnumIndex() != 1 {
			t.Errorf("Fill: Expected String 'FILL', EnumIndex 1; Got '%s', %d", typ.String(), typ.EnumIndex())
		}

		typ = Parity
		if typ.String() != "PARITY" || typ.EnumIndex() != 2 {
			t.Errorf("Parity: Expected String 'PARITY', EnumIndex 2; Got '%s', %d", typ.String(), typ.EnumIndex())
		}
	})

	// Test UnmarshalJSON()
	t.Run("UnmarshalJSON", func(t *testing.T) {
		var typ Type

		// Success case
		err := typ.UnmarshalJSON([]byte(`"FILL"`))
		if err != nil {
			t.Fatalf("UnmarshalJSON for 'FILL' failed: %v", err)
		}
		if typ != Fill {
			t.Errorf("UnmarshalJSON for 'FILL': Expected Fill, Got %v", typ)
		}

		// Error case: invalid string
		err = typ.UnmarshalJSON([]byte(`"INVALID"`))
		if err == nil {
			t.Fatalf("UnmarshalJSON for 'INVALID' expected an error, got nil")
		}
		if !errors.Is(err, ErrInvalidPaddingType) {
			t.Errorf("UnmarshalJSON for 'INVALID': Expected ErrInvalidPaddingType, Got %v", err)
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

		typ = Fill
		if !typ.IsValid() {
			t.Errorf("Fill: Expected IsValid true, Got false")
		}

		typ = Type(99) // Invalid type
		if typ.IsValid() {
			t.Errorf("Invalid Type: Expected IsValid false, Got true")
		}
	})
}

func TestPositionMethods(t *testing.T) {
	// Test String() and EnumIndex()
	t.Run("String and EnumIndex", func(t *testing.T) {
		var pos Position

		pos = Right
		if pos.String() != "RIGHT" || pos.EnumIndex() != 0 {
			t.Errorf("Right: Expected String 'RIGHT', EnumIndex 0; Got '%s', %d", pos.String(), pos.EnumIndex())
		}

		pos = Left
		if pos.String() != "LEFT" || pos.EnumIndex() != 1 {
			t.Errorf("Left: Expected String 'LEFT', EnumIndex 1; Got '%s', %d", pos.String(), pos.EnumIndex())
		}
	})

	// Test UnmarshalJSON()
	t.Run("UnmarshalJSON", func(t *testing.T) {
		var pos Position

		// Success case
		err := pos.UnmarshalJSON([]byte(`"LEFT"`))
		if err != nil {
			t.Fatalf("UnmarshalJSON for 'LEFT' failed: %v", err)
		}
		if pos != Left {
			t.Errorf("UnmarshalJSON for 'LEFT': Expected Left, Got %v", pos)
		}

		// Error case: invalid string
		err = pos.UnmarshalJSON([]byte(`"INVALID"`))
		if err == nil {
			t.Fatalf("UnmarshalJSON for 'INVALID' expected an error, got nil")
		}
		if !errors.Is(err, ErrInvalidPaddingPosition) {
			t.Errorf("UnmarshalJSON for 'INVALID': Expected ErrInvalidPaddingPosition, Got %v", err)
		}

		// Error case: non-string input
		err = pos.UnmarshalJSON([]byte(`123`))
		if err == nil {
			t.Fatalf("UnmarshalJSON for non-string expected an error, got nil")
		}
	})

	// Test IsValid()
	t.Run("IsValid", func(t *testing.T) {
		var pos Position

		pos = Left
		if !pos.IsValid() {
			t.Errorf("Left: Expected IsValid true, Got false")
		}

		pos = Position(99) // Invalid type
		if pos.IsValid() {
			t.Errorf("Invalid Position: Expected IsValid false, Got true")
		}
	})
}
