package bitmap

import (
	"fmt"
	"github.com/tomasdemarco/iso8583/packager/field"
	"github.com/tomasdemarco/iso8583/utils"
	"sort"
)

// Pack encodes a slice of field numbers (strings) into a bitmap byte slice.
// It determines if a primary (8 bytes) or secondary (16 bytes) bitmap is needed.
func Pack(fieldNumbers []int) (*utils.BitSet, error) {

	sort.Ints(fieldNumbers)

	numBits := 64
	hasSecondaryBitmap := false

	for _, num := range fieldNumbers {
		if num > 64 {
			hasSecondaryBitmap = true
			break
		}
	}

	if hasSecondaryBitmap {
		numBits = 128
	}

	bitmap := utils.NewBitSet(numBits, 128)

	for _, num := range fieldNumbers {
		bitmap.Set(num - 1)
	}

	if hasSecondaryBitmap {
		bitmap.Set(0)
	}

	return bitmap, nil
}

// Unpack decodes a bitmap byte slice into a slice of field numbers (strings).
// It determines if a primary (8 bytes) or secondary (16 bytes) bitmap is present.
func Unpack(field field.Packager, b []byte, offset int) (*utils.BitSet, int, error) {

	bitmapVal, length, err := field.Unpack(b, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("unpack primary bitmap: %w", err)
	}
	offset += length

	bMap, err := utils.Byte2BitSet(utils.Hex2Byte(bitmapVal))
	if err != nil {
		return nil, 0, fmt.Errorf("could not get bitmap: %w", err)
	}

	if bMap.Get(1) {
		bitmap2Val, length2, err := field.Unpack(b, offset)
		if err != nil {
			return nil, 0, fmt.Errorf("unpack secondary bitmap: %w", err)
		}
		bitmapVal += bitmap2Val
		length += length2

		bMap2, err := utils.Byte2BitSet(utils.Hex2Byte(bitmap2Val))
		if err != nil {
			return nil, 0, fmt.Errorf("could not get secondary bitmap: %w", err)
		}

		bMap = bMap.Concatenate(bMap2)
	}

	return bMap, length, nil
}
