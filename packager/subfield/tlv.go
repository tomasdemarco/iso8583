package subfield

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
)

type TLVPackager struct {
	left bool
	char string
}

func (f *TLVPackager) Pack(data map[string]string) (string, error) {
	var packedBytes bytes.Buffer

	// Get keys and sort them to ensure consistent order
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, tagStr := range keys {
		valueStr := data[tagStr]

		tagBytes, err := hex.DecodeString(tagStr)
		if err != nil {
			return "", fmt.Errorf("tag '%s' no es un hexadecimal válido: %w", tagStr, err)
		}

		valueBytes, err := hex.DecodeString(valueStr)
		if err != nil {
			return "", fmt.Errorf("value '%s' para tag '%s' no es un hexadecimal válido: %w", valueStr, tagStr, err)
		}

		length := len(valueBytes)

		if length > 255 {
			return "", fmt.Errorf("la longitud del valor para el tag '%s' (%d bytes) excede el máximo de 255 bytes permitidos para un campo Length de 1 byte", tagStr, length)
		}

		packedBytes.Write(tagBytes)
		packedBytes.WriteByte(byte(length))
		packedBytes.Write(valueBytes)
	}

	return hex.EncodeToString(packedBytes.Bytes()), nil
}

// Unpack convierte una cadena TLV a un map[string]string (tag-value).
func (f *TLVPackager) Unpack(tlvString string) (map[string]string, error) {
	unpackedData := make(map[string]string)

	// Primero, decodificar la cadena hexadecimal de entrada a bytes crudos.
	buf, err := hex.DecodeString(tlvString)
	if err != nil {
		return nil, fmt.Errorf("la cadena TLV de entrada no es un hexadecimal válido: %w", err)
	}

	offset := 0

	for offset < len(buf) {
		// Parse Tag
		var tagBytes []byte
		var tagLen int

		// Determinar la longitud del tag (1 o 2 bytes para EMV)
		if (buf[offset] & 0x1F) == 0x1F { // Si los últimos 5 bits del primer byte son 1s, es un tag de 2 bytes
			if offset+2 > len(buf) {
				return nil, errors.New("cadena TLV incompleta o mal formada: longitud insuficiente para Tag de 2 bytes")
			}
			tagBytes = buf[offset : offset+2]
			tagLen = 2
		} else {
			if offset+1 > len(buf) {
				return nil, errors.New("cadena TLV incompleta o mal formada: longitud insuficiente para Tag de 1 byte")
			}
			tagBytes = buf[offset : offset+1]
			tagLen = 1
		}
		tagStr := hex.EncodeToString(tagBytes)
		offset += tagLen

		// Parse Length (asumiendo 1-byte length como lo produce Pack)
		if offset+1 > len(buf) {
			return nil, errors.New("cadena TLV incompleta o mal formada: longitud insuficiente para Length")
		}
		length := int(buf[offset])
		offset++

		// Parse value
		if offset+length > len(buf) {
			return nil, fmt.Errorf("cadena TLV incompleta: el valor para el tag '%s' (longitud declarada %d) excede el tamaño restante de la cadena", tagStr, length)
		}
		valueBytes := buf[offset : offset+length]
		valueStr := hex.EncodeToString(valueBytes)
		offset += length

		unpackedData[tagStr] = valueStr
	}

	return unpackedData, nil
}
