package bitmap

import (
	"errors"
	"fmt"
	"github.com/tomasdemarco/iso8583/encoding"
	"github.com/tomasdemarco/iso8583/field"
	"github.com/tomasdemarco/iso8583/packager"
	"github.com/tomasdemarco/iso8583/utils"
	"sort"
	"strconv"
)

func Unpack(field1 packager.Field, position int, messageRaw []byte) (int, []string, error) {
	numberBitmaps := 1
	var lengthBitmap int

	var bitmapRaw string
	if field1.Encoding == encoding.Ascii {
		bitmapFirstChar, err := encoding.AsciiDecode(messageRaw[position : position+1])
		if err != nil {
			return 0, nil, err
		}

		validSecondBitmap, err := strconv.ParseInt(bitmapFirstChar, 16, 10)
		if err != nil {
			return 0, nil, err
		}

		if validSecondBitmap > 7 {
			numberBitmaps++
		}

		lengthBitmap = 16 * numberBitmaps

		if len(messageRaw) < position+lengthBitmap {
			return 0, nil, errors.New("index out of range when trying to get bitmap")
		}

		bitmapRaw, err = encoding.AsciiDecode(messageRaw[position : position+lengthBitmap])
		if err != nil {
			return 0, nil, err
		}
	} else if field1.Encoding == encoding.Ebcdic {
		bitmapFirstChar := encoding.EbcdicDecode(messageRaw[position : position+1])

		validSecondBitmap, err := strconv.ParseInt(bitmapFirstChar, 16, 10)
		if err != nil {
			return 0, nil, err
		}

		if validSecondBitmap > 7 {
			numberBitmaps++
		}

		lengthBitmap = 16 * numberBitmaps

		if len(messageRaw) < position+lengthBitmap {
			return 0, nil, errors.New("index out of range when trying to get bitmap")
		}

		bitmapRaw = encoding.EbcdicDecode(messageRaw[position : position+lengthBitmap])
	} else {
		validSecondBitmap, err := strconv.ParseInt(fmt.Sprintf("%x", messageRaw[position:position+1])[:1], 16, 10)
		if err != nil {
			return 0, nil, err
		}

		if validSecondBitmap > 7 {
			numberBitmaps++
		}

		lengthBitmap = 8 * numberBitmaps

		if len(messageRaw) < position+lengthBitmap {
			return 0, nil, errors.New("index out of range when trying to get bitmap")
		}

		bitmapRaw = fmt.Sprintf("%x", messageRaw[position:position+lengthBitmap])
	}

	sliceBitmap, err := encoding.BitmapDecode(utils.Hex2Byte(bitmapRaw))
	if err != nil {
		return 0, nil, err
	}

	return lengthBitmap, sliceBitmap, nil
}

func Pack(fields map[string]field.Field) ([]string, []byte, error) {
	sliceBitmap := make([]string, 0)

	for k := range fields {
		str := fmt.Sprintf("%03s", k)
		sliceBitmap = append(sliceBitmap, str)
	}

	sort.Strings(sliceBitmap)

	bitmap, err := encoding.BitmapEncode(sliceBitmap)
	if err != nil {
		return nil, nil, err
	}

	return sliceBitmap, bitmap, nil
}
