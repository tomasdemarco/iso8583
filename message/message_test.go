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
		initialFields map[string]string
	}{
		{
			name: "Mensaje simple con campos fijos",
			initialFields: map[string]string{
				"000": "0200",
				"011": "123456",
			},
		},
		{
			name: "Mensaje con campo de longitud variable (LLVAR)",
			initialFields: map[string]string{
				"000": "0800",
				"011": "654321",
				"032": "98765",
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
			fmt.Printf("%x\n\n", packedData)
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
					t.Errorf("El campo '%s' esperado no se encontró en el mensaje desempaquetado", id)
					continue
				}

				// Comparamos los valores, pero ignoramos el padding que se pudo haber añadido.
				if !strings.HasSuffix(finalValue, initialValue) {
					t.Errorf("El valor del campo '%s' no coincide. Esperado (sufijo): '%s', Recibido: '%s'", id, initialValue, finalValue)
				}
			}

			// 2. Comprobar que el campo del bitmap (001) fue creado.
			if _, err := msgToUnpack.GetField("001"); err != nil {
				t.Errorf("El campo del bitmap (001) no se encontró en el mensaje desempaquetado")
			}

			t.Logf("Prueba '%s' completada con éxito.", tc.name)
			t.Logf("Datos empaquetados (hex): %x", packedData)
			t.Logf("Campos desempaquetados: %v", msgToUnpack.Fields)
		})
	}
}
