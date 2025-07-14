package bitmap

import (
	"fmt"
	"github.com/tomasdemarco/iso8583/packager/field"
	"github.com/tomasdemarco/iso8583/utils"
)

// Unpack decodes a bitmap byte slice into a slice of field numbers (strings).
// It determines if a primary (8 bytes) or secondary (16 bytes) bitmap is present.
func Unpack(field field.Packager, b []byte, offset int) (*utils.BitSet, int, error) {

	bitmapVal, length, err := field.Unpack(b, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("%w: %w", ErrUnpackPrimaryBitmap, err)
	}
	offset += length

	bMap, err := utils.Byte2BitSet(utils.Hex2Byte(bitmapVal))
	if err != nil {
		return nil, 0, fmt.Errorf("%w: %w", ErrInvalidPrimaryBitmap, err)
	}

	if bMap.Get(1) {
		bitmap2Val, length2, err := field.Unpack(b, offset)
		if err != nil {
			return nil, 0, fmt.Errorf("%w: %w", ErrUnpackSecondaryBitmap, err)
		}
		bitmapVal += bitmap2Val
		length += length2

		bMap2, err := utils.Byte2BitSet(utils.Hex2Byte(bitmap2Val))
		if err != nil {
			return nil, 0, fmt.Errorf("%w: %w", ErrInvalidSecondaryBitmap, err)
		}

		bMap = bMap.Concatenate(bMap2)
	}

	return bMap, length, nil
}
