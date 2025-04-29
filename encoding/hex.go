package encoding

import (
	"fmt"
	"github.com/tomasdemarco/iso8583/utils"
	"strconv"
)

// HEX implements the Encoder interface for HEX encoding.
type HEX struct {
	length int
}

func NewHexEncoder() HEX {
	return HEX{}
}

func (e *HEX) Encode(src string) ([]byte, error) {
	intValue, err := strconv.ParseInt(src, 10, 16)
	if err != nil {
		return nil, err
	}

	stringValue := strconv.FormatInt(intValue, 16)
	if len(stringValue)%2 != 0 {
		stringValue = "0" + stringValue
	}

	return utils.Hex2Byte(stringValue), nil
}

func (e *HEX) Decode(src []byte) (string, error) {
	dstInt64, err := strconv.ParseInt(fmt.Sprintf("%x", src[:e.length]), 16, 32)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v", dstInt64), nil
}

func (e *HEX) SetLength(length int) {
	e.length = length
}
