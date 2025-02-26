package encoding

import (
	"fmt"
	"strconv"
)

func BitmapDecode(value string) ([]string, error) {

	primaryBitmap, err := strconv.ParseUint(value[:16], 16, 64)
	if err != nil {
		return nil, err
	}

	bitmapBinary := fmt.Sprintf("%064b", primaryBitmap)

	if len(value) == 32 {
		secondaryBitmap, err := strconv.ParseUint(value[16:], 16, 64)
		if err != nil {
			return nil, err
		}

		bitmapBinary += fmt.Sprintf("%064b", secondaryBitmap)
	}

	sliceBitmap := make([]string, 0)

	for i := 1; i <= len(bitmapBinary); i++ {
		if bitmapBinary[i-1:i] == "1" {
			str := fmt.Sprintf("%03d", i)
			sliceBitmap = append(sliceBitmap, str)
		}
	}

	return sliceBitmap, nil
}

func BitmapEncode(value []string) (string, error) {
	var bitmap string
	bitmapArray := make(map[int]string)

	numberBits := 64

	for _, i := range value {
		val, err := strconv.Atoi(i)
		if err != nil {
			return bitmap, err
		}

		bitmapArray[val] = "1"

		if val > 64 {
			numberBits = 128
			bitmapArray[1] = "1"
		}
	}

	for i := 1; i <= numberBits; i++ {
		if _, ok := bitmapArray[i]; ok {
			bitmap += "1"
		} else {
			bitmap += "0"
		}
	}

	primaryBitmap, err := strconv.ParseUint(bitmap[:64], 2, 64)
	if err != nil {
		return bitmap, err
	}
	result := fmt.Sprintf("%016x", primaryBitmap)

	if len(bitmap) > 64 {
		secondary, err := strconv.ParseUint(bitmap[64:], 2, 64)
		if err != nil {
			return bitmap, err
		}
		result += fmt.Sprintf("%016x", secondary)
	}

	return result, nil
}
