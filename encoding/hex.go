package encoding

import (
	"fmt"
	"strconv"
)

func HexDecode(value string) (string, error) {
	stringValue, err := strconv.ParseInt(value, 16, 32)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", stringValue), nil
}

func HexEncode(value string) (string, error) {
	intValue, err := strconv.ParseInt(value, 10, 16)
	if err != nil {
		return "", err
	}
	stringValue := strconv.FormatInt(intValue, 16)
	return stringValue, nil
}
