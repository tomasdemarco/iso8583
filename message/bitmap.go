package message

import (
	"errors"
	"fmt"
	"strconv"
)

func (m *Message) UnpackBitmap(positionInitial int, messageRaw string) (int, []string, error) {
	numberBitmaps := 1
	validSecondBitmap, err := strconv.ParseInt(messageRaw[positionInitial:positionInitial+1], 10, 16)
	if validSecondBitmap > 8 || err != nil {
		numberBitmaps++
	}

	if len(messageRaw) < m.Packager.Fields["001"].Length*numberBitmaps {
		err = errors.New("index out of range when trying to get bitmap")
		return 0, nil, err
	}

	bitmapBinary := ""
	increase := 0
	for i := 0; i < numberBitmaps; i++ {
		bitmap := messageRaw[positionInitial+increase : m.Packager.Fields["001"].Length+positionInitial+increase]
		firstBitmap, err := strconv.ParseUint(bitmap[:8], 16, 32)
		if err != nil {
			return 0, nil, err
		}
		firstBitmapAux, err := strconv.ParseUint(bitmap[8:16], 16, 32)
		if err != nil {
			return 0, nil, err
		}
		bitmapBinary += fmt.Sprintf("%032b%032b", firstBitmap, firstBitmapAux)
		increase += 16
	}
	sliceBitmap := make([]string, 0)

	for i := 1; i <= len(bitmapBinary); i++ {
		if bitmapBinary[i-1:i] == "1" {
			str := fmt.Sprintf("%03d", i)
			sliceBitmap = append(sliceBitmap, str)
		}
	}
	return numberBitmaps, sliceBitmap, nil
}

func (m *Message) PackBitmap() string {
	numberBitmaps := 1
	numberBits := 64
	bitmap := "0"

	var bitmapArray = []int{}

	for _, i := range m.Bitmap {
		j, err := strconv.Atoi(i)
		if err != nil {
			panic(err)
		}
		bitmapArray = append(bitmapArray, j)
	}

	for _, num := range bitmapArray {
		if num > 64 {
			numberBits = 128
			numberBitmaps = 2
			bitmap = "1"
		}
	}

	for i := 2; i <= numberBits; i++ {
		if _, ok := m.FieldAndSubFields[fmt.Sprintf("%03d", i)]; ok {
			bitmap += "1"
		} else {
			bitmap += "0"
		}
	}
	var result string
	aux := 0
	for i := 0; i < 16*numberBitmaps; i++ {
		bitmap1, _ := strconv.ParseInt(bitmap[aux:4+aux], 2, 64)
		result += fmt.Sprintf("%x", bitmap1)
		aux += 4
	}
	return result
}
