package message

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/tomasdemarco/iso8583/packager"
	"github.com/tomasdemarco/iso8583/utils"
)

// BitmapCustomField es un CustomPacker de ejemplo para un campo que contiene subcampos
// cuya presencia es definida por un bitmap interno.
type BitmapCustomField struct {
	SubfieldValues map[int]string     // Valores de los subcampos (ID -> Valor)
	SubPackager    *packager.Packager // FieldPackager que define la estructura de los subcampos
	InternalBitmap *utils.BitSet      // El bitmap que controla la presencia de subcampos
}

// GetValue helper para BitmapCustomField
func (f *BitmapCustomField) GetValue(id int) (string, bool) {
	val, ok := f.SubfieldValues[id]
	return val, ok && f.InternalBitmap != nil && f.InternalBitmap.Get(id) // Solo si está activo en el bitmap
}

// SetValue helper para BitmapCustomField
func (f *BitmapCustomField) SetValue(id int, value string) {
	if f.SubfieldValues == nil {
		f.SubfieldValues = make(map[int]string)
	}
	f.SubfieldValues[id] = value

	if f.InternalBitmap == nil {
		// Asume un tamaño inicial de bitmap adecuado, por ejemplo 64/128 bits
		f.InternalBitmap = utils.NewBitSet(f.SubPackager.Fields[0].Length()*8, 128)
	}
	f.InternalBitmap.Set(id) // Marca el subcampo como presente
}

// Pack implementa CustomPacker.Pack() para BitmapCustomField.
func (f *BitmapCustomField) Pack() (string, error) {
	if f.SubPackager == nil {
		return "", fmt.Errorf("BitmapCustomField: SubPackager no inicializado para empaquetar")
	}
	if f.InternalBitmap == nil || f.InternalBitmap.GetSize() == 0 {
		return "", nil // Si no hay bitmap o está vacío, no hay nada que empaquetar
	}

	var packedData bytes.Buffer

	// Empaquetar el bitmap interno. Asumimos que el SubPackager tiene un campo
	// definido para el bitmap, usando .ToString() del BitSet.
	if bitmapFldPkg, ok := f.SubPackager.Fields[0]; ok {
		packedInternalBitmap, _, err := bitmapFldPkg.Pack(f.InternalBitmap.ToString())
		if err != nil {
			return "", fmt.Errorf("error al empaquetar el bitmap interno: %w", err)
		}
		packedData.Write(packedInternalBitmap)
	} else {
		// Alternativa: Si no hay un campo 1 definido, simplemente empaqueta raw bytes del bitmap.
		// Esto dependerá mucho de cómo esté definido tu packager de subcampos.
		// packedData.Write(f.InternalBitmap.ToBytes())
		return "", fmt.Errorf("BitmapCustomField: SubPackager no tiene definición para campo bitmap (ID 0)")
	}

	// Luego, empaquetamos los subcampos que están activos en el InternalBitmap
	for _, id := range f.InternalBitmap.GetSliceString() {
		if id == 0 { // El campo 0 (bitmap) ya ha sido manejado
			continue
		}

		if val, ok := f.SubfieldValues[id]; ok {
			if subFldPkg, ok := f.SubPackager.Fields[id]; ok {
				packedSubfield, _, err := subFldPkg.Pack(val)
				if err != nil {
					return "", fmt.Errorf("error al empaquetar subcampo %d: %w", id, err)
				}
				packedData.Write(packedSubfield)
			} else {
				return "", fmt.Errorf("subcampo %d no definido en SubPackager", id)
			}
		}
	}

	return hex.EncodeToString(packedData.Bytes()), nil
}

// Unpack implementa CustomPacker.Unpack() para BitmapCustomField.
func (f *BitmapCustomField) Unpack(fieldData string) error {
	if f.SubPackager == nil {
		return fmt.Errorf("BitmapCustomField: SubPackager no inicializado para desempaquetar")
	}

	rawBytes, err := hex.DecodeString(fieldData)
	if err != nil {
		return fmt.Errorf("BitmapCustomField: error al decodificar string hexadecimal: %w", err)
	}

	// Creamos un "sub-mensaje" para desempaquetar el bitmap interno y los subcampos
	subMessage := NewMessage(f.SubPackager)
	// Solo desempaquetamos si rawBytes no está vacío
	if len(rawBytes) > 0 {
		if err := subMessage.Unpack(rawBytes); err != nil {
			return fmt.Errorf("error al desempaquetar subcampos: %w", err)
		}
	} else {
		// Si no hay datos binarios, asumimos que no hay subcampos para desempaquetar
		f.InternalBitmap = utils.NewBitSet(f.SubPackager.Fields[0].Length()*8, 128)
		f.SubfieldValues = make(map[int]string)
		return nil
	}

	// El bitmap interno es el Bitmap desempaquetado del sub-mensaje
	f.InternalBitmap = subMessage.Bitmap
	f.SubfieldValues = make(map[int]string)

	// Ahora iteramos los campos activos en el InternalBitmap para obtener sus valores
	for _, id := range f.InternalBitmap.GetSliceString() {
		// MTI (0) y el bitmap en sí (1) no son "valores de subcampo" que almacenemos directamente aquí
		if id == 0 {
			continue
		}
		if subField, err := subMessage.Field(id).String(); err != nil {
			f.SubfieldValues[id] = subField
		}
	}

	return nil
}

func (f *BitmapCustomField) Log() (interface{}, error) {
	return f.SubfieldValues, nil
}
