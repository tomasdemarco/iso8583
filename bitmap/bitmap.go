package bitmap

import (
	"fmt"
	"github.com/tomasdemarco/iso8583/encoding"
	pkgField "github.com/tomasdemarco/iso8583/packager/field"
	"github.com/tomasdemarco/iso8583/utils"
	"sort"
)

func Unpack(pkgField pkgField.Field, messageRaw []byte, position int) (*int, []string, error) {

	length := pkgField.Length
	if pkgField.Encoding == encoding.Binary {
		length = length / 2
	}

	value, err := encoding.Unpack(pkgField.Encoding, messageRaw[position:position+length])
	if err != nil {
		return nil, nil, err
	}

	sliceBitmap, err := encoding.BitmapDecode(utils.Hex2Byte(value), 1)
	if err != nil {
		return nil, nil, err
	}

	if sliceBitmap[0] == "001" {
		secLength := pkgField.Length
		if pkgField.Encoding == encoding.Binary {
			secLength -= pkgField.Length / 2
		}

		value, err = encoding.Unpack(pkgField.Encoding, messageRaw[position+length:position+length+secLength])
		if err != nil {
			return nil, nil, err
		}

		length += secLength

		sliceExtBitmap, err := encoding.BitmapDecode(utils.Hex2Byte(value), 65)
		if err != nil {
			return nil, nil, err
		}

		sliceBitmap = append(sliceBitmap, sliceExtBitmap...)
	}

	return &length, sliceBitmap, nil
}

func Pack(fields map[string]string) ([]string, []byte, error) {
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
