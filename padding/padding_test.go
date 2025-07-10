package padding

import (
	"testing"

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
			fieldLength:      6, // Longitud total del campo (par)
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
