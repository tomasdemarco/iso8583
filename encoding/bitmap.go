package encoding

import (
	"encoding/binary"
	"fmt"
	"strconv"
)

// BitmapDecode decodifica un slice de bytes de bitmap en un slice de strings de campos.
func BitmapDecode(bitmapBytes []byte, initBit int) ([]string, error) {
	if len(bitmapBytes) < 8 {
		return nil, fmt.Errorf("se requieren al menos 8 bytes para el bitmap primario, se obtuvieron %d", len(bitmapBytes))
	}

	primaryBitmap := binary.BigEndian.Uint64(bitmapBytes[:8])
	bitmapBinary := fmt.Sprintf("%064b", primaryBitmap)

	// El bit 1 (el más significativo del primer byte) indica la presencia de un bitmap secundario.
	// Si el bit 1 está activo, el valor del primer byte será >= 0x80.
	if (bitmapBytes[0] & 0x80) != 0 { // Comprobar si el bit 1 está activo
		if len(bitmapBytes) < 16 {
			return nil, fmt.Errorf("se requieren 16 bytes para el bitmap secundario, se obtuvieron %d", len(bitmapBytes))
		}
		secondaryBitmap := binary.BigEndian.Uint64(bitmapBytes[8:16])
		bitmapBinary += fmt.Sprintf("%064b", secondaryBitmap)
	}

	sliceBitmap := make([]string, 0)

	for i, bit := range bitmapBinary {
		if bit == '1' {
			str := fmt.Sprintf("%03d", i+initBit)
			sliceBitmap = append(sliceBitmap, str)
		}
	}

	return sliceBitmap, nil
}

// BitmapEncode codifica un slice de strings de campos en un slice de bytes de bitmap.
func BitmapEncode(fieldNumbers []string) ([]byte, error) {
	// Determinar si se necesita un bitmap secundario (si hay campos > 64)
	needsSecondary := false
	for _, fn := range fieldNumbers {
		num, err := strconv.Atoi(fn)
		if err != nil {
			return nil, fmt.Errorf("número de campo inválido: %s", fn)
		}
		if num > 64 {
			needsSecondary = true
			break
		}
	}

	// Crear un slice de booleanos para representar los bits
	numBits := 64
	if needsSecondary {
		numBits = 128
	}
	bits := make([]bool, numBits+1) // +1 para usar índices 1-basados

	// Marcar los bits correspondientes a los números de campo
	for _, fn := range fieldNumbers {
		num, _ := strconv.Atoi(fn)
		if num > 0 && num <= numBits {
			bits[num] = true
		}
	}

	// Activar el bit 1 si hay un bitmap secundario
	// Esto se hace si needsSecondary es true, ya que el bit 1 indica la presencia del secundario.
	if needsSecondary {
		bits[1] = true
	}

	// Convertir el slice de bits a un string binario
	var binaryString string
	for i := 1; i <= numBits; i++ {
		if bits[i] {
			binaryString += "1"
		} else {
			binaryString += "0"
		}
	}

	// Convertir el string binario a bytes
	var result []byte
	primaryBinary := binaryString[:64]
	primaryUint, err := strconv.ParseUint(primaryBinary, 2, 64)
	if err != nil {
		return nil, fmt.Errorf("error al parsear bitmap primario: %w", err)
	}
	primaryBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(primaryBytes, primaryUint)
	result = append(result, primaryBytes...)

	if needsSecondary {
		secondaryBinary := binaryString[64:]
		secondaryUint, err := strconv.ParseUint(secondaryBinary, 2, 64)
		if err != nil {
			return nil, fmt.Errorf("error al parsear bitmap secundario: %w", err)
		}
		secondaryBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(secondaryBytes, secondaryUint)
		result = append(result, secondaryBytes...)
	}

	return result, nil
}
