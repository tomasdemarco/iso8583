package message

import (
	"errors"
	"fmt"
	"gitlab.com/g6604/adquirencia/desarrollo/golang_package/iso8583/encoding"
	"sort"
	"strconv"
)

func (m *Message) UnpackBitmap(position int, messageRaw string) (int, []string, error) {
	numberBitmaps := 1
	var lengthBitmap int

	var bitmapRaw string
	if m.Packager.Fields["001"].Encoding == encoding.Ascii {
		bitmapFirstChar, err := encoding.AsciiDecode(messageRaw[position : position+2])
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

		lengthBitmap = 32 * numberBitmaps

		if len(messageRaw) < position+lengthBitmap {
			return 0, nil, errors.New("index out of range when trying to get bitmap")
		}

		bitmapRaw, err = encoding.AsciiDecode(messageRaw[position : position+lengthBitmap])
		if err != nil {
			return 0, nil, err
		}
	} else if m.Packager.Fields["001"].Encoding == encoding.Ebcdic {
		bitmapFirstChar, err := encoding.EbcdicDecode(messageRaw[position : position+2])
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

		lengthBitmap = 32 * numberBitmaps

		if len(messageRaw) < position+lengthBitmap {
			return 0, nil, errors.New("index out of range when trying to get bitmap")
		}

		bitmapRaw, err = encoding.EbcdicDecode(messageRaw[position : position+lengthBitmap])
		if err != nil {
			return 0, nil, err
		}
	} else {
		validSecondBitmap, err := strconv.ParseInt(messageRaw[position:position+1], 16, 10)
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

		bitmapRaw = messageRaw[position : position+lengthBitmap]
	}

	sliceBitmap, err := encoding.BitmapDecode(bitmapRaw)
	if err != nil {
		return 0, nil, err
	}

	return lengthBitmap, sliceBitmap, nil
}

func (m *Message) PackBitmap() (string, error) {
	sliceBitmap := make([]string, 0)

	for k := range m.Fields {
		str := fmt.Sprintf("%03s", k)
		sliceBitmap = append(sliceBitmap, str)
	}

	sort.Strings(sliceBitmap)

	m.Bitmap = sliceBitmap

	bitmap, err := encoding.BitmapEncode(m.Bitmap)
	if err != nil {
		return bitmap, err
	}

	return bitmap, nil
}
