package encoding

import (
	"fmt"
	"strconv"
)

func HexDecode(src []byte) (string, error) {
	dstInt64, err := strconv.ParseInt(fmt.Sprintf("%x", src), 16, 32)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", dstInt64), nil
}

func HexEncode(value string) ([]byte, error) {
	intValue, err := strconv.ParseInt(value, 10, 16)
	if err != nil {
		return nil, err
	}
	stringValue := strconv.FormatInt(intValue, 16)
	return []byte(stringValue), nil
}
