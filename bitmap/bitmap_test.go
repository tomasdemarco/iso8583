package bitmap

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/tomasdemarco/iso8583/encoding"
	pkgField "github.com/tomasdemarco/iso8583/packager/field"
)

func TestBitmapRoundTrip(t *testing.T) {
	testCases := []struct {
		name                string
		fields              map[string]string
		field001Encoder     encoding.Encoder // Encoder para el campo 001 (bitmap)
		expectedPackedBytes []byte
		expectedBitmapSlice []string
	}{
		{
			name: "Primary Bitmap - Fields 2, 3, 5 (BINARY)",
			fields: map[string]string{
				"002": "value2",
				"003": "value3",
				"005": "value5",
			},
			field001Encoder:     &encoding.BINARY{},
			expectedPackedBytes: []byte{0x68, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			expectedBitmapSlice: []string{"002", "003", "005"},
		},
		{
			name: "Primary + Secondary Bitmap - Fields 1, 65, 128 (BINARY)",
			fields: map[string]string{
				"001": "value1", // Campo 1 activa el bitmap secundario
				"065": "value65",
				"128": "value128",
			},
			field001Encoder:     &encoding.BINARY{},
			expectedPackedBytes: []byte{0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
			expectedBitmapSlice: []string{"001", "065", "128"},
		},
		{
			name: "Primary Bitmap - Fields 2, 3 (ASCII Encoded)",
			fields: map[string]string{
				"002": "value2",
				"003": "value3",
			},
			field001Encoder: &encoding.ASCII{},
			// Bitmap binario para 2,3 es 01100000... (0x60). Hex string "6000000000000000"
			// ASCII de "6000000000000000" es 0x36, 0x30, 0x30, ...
			expectedPackedBytes: []byte{0x36, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30},
			expectedBitmapSlice: []string{"002", "003"},
		},
		{
			name: "Primary Bitmap - Fields 2, 3 (EBCDIC Encoded)",
			fields: map[string]string{
				"002": "value2",
				"003": "value3",
			},
			field001Encoder: &encoding.EBCDIC{},
			// Bitmap binario para 2,3 es 01100000... (0x60). Hex string "6000000000000000"
			// EBCDIC de "6000000000000000" es 0xF6, 0xF0, 0xF0, ...
			expectedPackedBytes: []byte{0xF6, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0},
			expectedBitmapSlice: []string{"002", "003"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// --- Test Pack ---
			packedBitmapSlice, packedBitmapBytes, err := Pack(tc.fields)
			if err != nil {
				t.Fatalf("Pack() falló: %v", err)
			}

			// Verificar el slice de bitmap generado
			if !reflect.DeepEqual(packedBitmapSlice, tc.expectedBitmapSlice) {
				t.Errorf("Pack() slice de bitmap incorrecto.\nEsperado: %v\nRecibido:  %v", tc.expectedBitmapSlice, packedBitmapSlice)
			}

			// Verificar los bytes del bitmap empaquetado (usando el encoder del caso de prueba)
			// Primero, codificamos el string hexadecimal del bitmap con el encoder del caso de prueba
			tc.field001Encoder.SetLength(len(packedBitmapBytes) * 2) // Longitud en caracteres hex
			encodedBitmapForComparison, err := tc.field001Encoder.Encode(fmt.Sprintf("%X", packedBitmapBytes))
			if err != nil {
				t.Fatalf("Error al codificar el bitmap para comparación: %v", err)
			}

			if !bytes.Equal(encodedBitmapForComparison, tc.expectedPackedBytes) {
				t.Errorf("Pack() bytes empaquetados incorrectos.\nEsperado: %x\nRecibido:  %x", tc.expectedPackedBytes, encodedBitmapForComparison)
			}

			// --- Test Unpack ---
			// Necesitamos un pkgField.Field para Unpack
			field001 := pkgField.Field{
				Length:   len(tc.expectedPackedBytes), // Longitud del bitmap en bytes
				Encoding: tc.field001Encoder,
			}

			_, unpackedBitmapSlice, err := Unpack(field001, tc.expectedPackedBytes, 0)
			if err != nil {
				t.Fatalf("Unpack() falló: %v", err)
			}

			// Verificar el slice de bitmap desempaquetado
			if !reflect.DeepEqual(unpackedBitmapSlice, tc.expectedBitmapSlice) {
				t.Errorf("Unpack() slice de bitmap incorrecto.\nEsperado: %v\nRecibido:  %v", tc.expectedBitmapSlice, unpackedBitmapSlice)
			}
		})
	}
}
