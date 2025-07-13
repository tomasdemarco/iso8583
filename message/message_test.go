package message

import (
	"fmt"
	"strings"
	"testing"

	"github.com/tomasdemarco/iso8583/packager"
)

func TestPackUnpackRoundTrip(t *testing.T) {
	pkg, err := packager.LoadFromJson(".", "test_packager.json")
	if err != nil {
		t.Fatalf("Error al cargar el packager de prueba: %v", err)
	}

	testCases := []struct {
		name          string
		initialFields map[int]string
	}{
		{
			name: "Mensaje con campo de longitud variable (LLVAR)",
			initialFields: map[int]string{
				0:  "0800",
				11: "654321",
				32: "98765",
			},
		},
		{
			name: "Mensaje simple con campos fijos",
			initialFields: map[int]string{
				0:  "0200",
				11: "123456",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// --- Fase de PACK ---
			msgToPack := NewMessage(pkg)
			for id, value := range tc.initialFields {
				msgToPack.SetField(id, value)
			}

			packedData, err := msgToPack.Pack()
			if err != nil {
				t.Fatalf("Pack() falló: %v", err)
			}
			fmt.Printf("Packed data: %x\n", packedData)
			// --- Fase de UNPACK ---
			msgToUnpack := NewMessage(pkg)
			err = msgToUnpack.Unpack(packedData)
			if err != nil {
				t.Fatalf("Unpack() falló: %v", err)
			}

			// --- Verificación ---

			// 1. Comprobar que todos los campos iniciales existen en el resultado
			for id, initialValue := range tc.initialFields {
				finalValue, err := msgToUnpack.GetField(id)
				if err != nil {
					t.Errorf("El campo '%d' esperado no se encontró en el mensaje desempaquetado", id)
					continue
				}

				// Comparamos los valores, pero ignoramos el padding que se pudo haber añadido.
				if !strings.HasSuffix(finalValue, initialValue) {
					t.Errorf("El valor del campo '%d' no coincide. Esperado (sufijo): '%s', Recibido: '%s'", id, initialValue, finalValue)
				}
			}

			// 2. Comprobar que el campo del bitmap (001) fue creado.
			if _, err := msgToUnpack.GetField(1); err != nil {
				t.Errorf("El campo del bitmap (1) no se encontró en el mensaje desempaquetado")
			}

			t.Logf("Prueba '%s' completada con éxito.", tc.name)
			t.Logf("Datos empaquetados (hex): %x", packedData)
			t.Logf("Campos desempaquetados: %v", msgToUnpack.Fields)
		})
	}
}
