package message

import (
	"errors"
	"fmt"
	"github.com/tomasdemarco/iso8583/encoding"
	"sort"
	"strconv"
)

func (m *Message) UnpackBitmap(positionInitial int, messageRaw string) (int, []string, error) {
	numberBitmaps := 1

	var bitmapRaw string
	if m.Packager.Fields["001"].Encoding == "ASCII" {
		bitmapFirstChar, err := encoding.AsciiDecode(messageRaw[positionInitial : positionInitial+2])
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

		bitmapRaw, err = encoding.AsciiDecode(messageRaw[positionInitial : positionInitial+(32*numberBitmaps)])
		if err != nil {
			return 0, nil, err
		}

		if len(bitmapRaw) < m.Packager.Fields["001"].Length*2*numberBitmaps {
			return 0, nil, errors.New("index out of range when trying to get bitmap")
		}
	} else if m.Packager.Fields["001"].Encoding == "EBCDIC" {
		bitmapFirstChar, err := encoding.EbcdicDecode(messageRaw[positionInitial : positionInitial+2])
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

		bitmapRaw, err = encoding.EbcdicDecode(messageRaw[positionInitial : positionInitial+(32*numberBitmaps)])
		if err != nil {
			return 0, nil, err
		}

		if len(bitmapRaw) < m.Packager.Fields["001"].Length*2*numberBitmaps {
			return 0, nil, errors.New("index out of range when trying to get bitmap")
		}
	} else {
		validSecondBitmap, err := strconv.ParseInt(messageRaw[positionInitial:positionInitial+1], 16, 10)
		if err != nil {
			return 0, nil, err
		}

		if validSecondBitmap > 7 {
			numberBitmaps++
		}

		bitmapRaw = messageRaw[positionInitial : positionInitial+(16*numberBitmaps)]

		if len(bitmapRaw) < m.Packager.Fields["001"].Length*numberBitmaps {
			return 0, nil, errors.New("index out of range when trying to get bitmap")
		}
	}

	sliceBitmap, err := encoding.BitmapDecode(bitmapRaw)
	if err != nil {
		return 0, nil, err
	}

	return numberBitmaps, sliceBitmap, nil
}

func (m *Message) PackBitmap() (string, error) {
	sliceBitmap := make([]string, 0)

	for k := range m.Fields {
		str := fmt.Sprintf("%03s", k)
		sliceBitmap = append(sliceBitmap, str)
	}

	sort.Strings(sliceBitmap)
	m.Bitmap = sliceBitmap
	fmt.Println(m.Bitmap)
	bitmap, err := encoding.BitmapEncode(m.Bitmap)
	if err != nil {
		return bitmap, err
	}

	return bitmap, nil
}
