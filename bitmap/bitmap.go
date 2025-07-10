// Package bitmap provides functionalities for encoding and decoding ISO 8583 bitmaps.
package bitmap

import (
	"fmt"
	"github.com/tomasdemarco/iso8583/encoding"
	pkgField "github.com/tomasdemarco/iso8583/packager/field"
	"github.com/tomasdemarco/iso8583/utils"
	"sort"
)

// Unpack decodes a bitmap from a raw message byte slice.
// It takes the field definition for the bitmap (pkgField), the raw message bytes (messageRaw),
// and the starting position of the bitmap in the message.
// It returns the total length of the bitmap in bytes, a slice of active field numbers (e.g., "002", "003"),
// and an error if unpacking fails.
func Unpack(pkgField pkgField.Field, messageRaw []byte, position int) (int, []string, error) {
	// The bitmap length is typically 8 bytes for primary, 16 for primary+secondary.
	// pkgField.Length here refers to the length of field 001 in the packager,
	// which should be 8 or 16.
	bitmapLength := pkgField.Length

	if len(messageRaw) < position+bitmapLength {
		return 0, nil, fmt.Errorf("insufficient data to unpack bitmap: expected %d bytes, got %d", bitmapLength, len(messageRaw)-position)
	}

	// Extract the encoded bitmap bytes.
	encodedBitmapBytes := messageRaw[position : position+bitmapLength]

	// Decode the encoded bytes to their string representation (hexadecimal).
	pkgField.Encoding.SetLength(bitmapLength)
	bitmapHexString, err := pkgField.Encoding.Decode(encodedBitmapBytes)
	if err != nil {
		return 0, nil, fmt.Errorf("error decoding encoded bitmap: %w", err)
	}

	// Convert the hexadecimal string to binary bytes for BitmapDecode.
	decodedBitmapBytes := utils.Hex2Byte(bitmapHexString)

	sliceBitmap, err := encoding.BitmapDecode(decodedBitmapBytes, 1)
	if err != nil {
		return 0, nil, fmt.Errorf("error decoding bitmap: %w", err)
	}

	return bitmapLength, sliceBitmap, nil
}

// Pack encodes a map of fields into a bitmap byte slice.
// It takes a map of field IDs to their values.
// It returns a sorted slice of active field numbers, the packed bitmap as a byte slice,
// and an error if packing fails.
func Pack(fields map[string]string) ([]string, []byte, error) {
	sliceBitmap := make([]string, 0)

	for k := range fields {
		str := fmt.Sprintf("%03s", k)
		sliceBitmap = append(sliceBitmap, str)
	}

	sort.Strings(sliceBitmap)

	bitmapBytes, err := encoding.BitmapEncode(sliceBitmap)
	if err != nil {
		return nil, nil, fmt.Errorf("error encoding bitmap: %w", err)
	}

	return sliceBitmap, bitmapBytes, nil
}
