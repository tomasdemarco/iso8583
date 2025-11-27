package message

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/tomasdemarco/iso8583/packager"
	"sort"
)

// TLVCustomField es un CustomPacker de ejemplo para un campo que contiene subcampos TLV.
// Almacena los subcampos como un map[int]string.
// Asume que el SubPackager está configurado para un empaquetado/desempaquetado secuencial.
// Si tu TLV es de un tipo fijo, puedes usar una struct en lugar del map.
type TLVCustomField struct {
	SubfieldValues map[int]string     // Valores de los subcampos (ID -> Valor)
	SubPackager    *packager.Packager // FieldPackager que define la estructura de los subcampos TLV
}

// GetValue helper para TLVCustomField
func (f *TLVCustomField) GetValue(id int) (string, bool) {
	val, ok := f.SubfieldValues[id]
	return val, ok
}

// SetValue helper para TLVCustomField
func (f *TLVCustomField) SetValue(id int, value string) {
	if f.SubfieldValues == nil {
		f.SubfieldValues = make(map[int]string)
	}
	f.SubfieldValues[id] = value
}

// Pack implementa CustomPacker.Pack() para TLVCustomField.
func (f *TLVCustomField) Pack() (string, error) {
	if f.SubPackager == nil {
		return "", fmt.Errorf("TLVCustomField: SubPackager no inicializado para empaquetar")
	}
	if f.SubfieldValues == nil || len(f.SubfieldValues) == 0 {
		return "", nil // Si no hay valores, devuelve string vacío
	}

	var packedBytesBuffer bytes.Buffer
	// El orden de empaquetado de los subcampos es crucial para TLV.
	// Aquí asumimos que los IDs de tus subcampos TLV son los números de campo ISO
	// y que se empaquetan en orden creciente.
	for _, id := range f.getSortedSubfieldKeys() {
		if val, ok := f.SubfieldValues[id]; ok {
			if subFldPkg, pkgOk := f.SubPackager.Fields[id]; pkgOk {
				packedSubfield, _, err := subFldPkg.Pack(val)
				if err != nil {
					return "", fmt.Errorf("error al empaquetar subcampo TLV %d: %w", id, err)
				}
				packedBytesBuffer.Write(packedSubfield)
			} else {
				return "", fmt.Errorf("subcampo TLV %d no definido en SubPackager", id)
			}
		}
	}

	// Devuelve la representación en hexadecimal del campo TLV completo.
	// Nota: Si tus campos TLV usan un formato específico de Tag+Length+Value,
	// deberás implementar esa lógica de composición aquí antes de la codificación hex.
	return hex.EncodeToString(packedBytesBuffer.Bytes()), nil
}

// Unpack implementa CustomPacker.Unpack() para TLVCustomField.
func (f *TLVCustomField) Unpack(data string) error {
	if f.SubPackager == nil {
		return fmt.Errorf("TLVCustomField: SubPackager no inicializado para desempaquetar")
	}

	rawBytes, err := hex.DecodeString(data)
	if err != nil {
		return fmt.Errorf("TLVCustomField: error al decodificar string hexadecimal: %w", err)
	}

	// Creamos un "sub-mensaje" usando el SubPackager para desempaquetar los subcampos.
	// Esto es útil si los subcampos se comportan como un mini-mensaje ISO.
	subMessage := NewMessage(f.SubPackager)
	if err := subMessage.Unpack(rawBytes); err != nil {
		return fmt.Errorf("error al desempaquetar subcampos TLV: %w", err)
	}

	f.SubfieldValues = make(map[int]string)
	// Cargamos los valores desempaquetados de nuevo a nuestro mapa
	for _, id := range subMessage.Bitmap.GetSliceString() {
		if subField, err := subMessage.Field(id).String(); err != nil { // Asegurarse de que el campo exista
			f.SubfieldValues[id] = subField
		}
	}
	return nil
}

// getSortedSubfieldKeys es un helper interno para asegurar un orden consistente al empaquetar.
func (f *TLVCustomField) getSortedSubfieldKeys() []int {
	keys := make([]int, 0, len(f.SubfieldValues))
	for id := range f.SubfieldValues {
		keys = append(keys, id)
	}
	sort.Ints(keys) // Ordena las keys para un empaquetado consistente
	return keys
}

func (f *TLVCustomField) Log() (interface{}, error) {
	return f.SubfieldValues, nil
}
