package encoding

import (
	"fmt"
	"strconv"
)

func HexDecode(value string) string {
	stringValue, err := strconv.ParseInt(value, 16, 32)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%v", stringValue)
}

func HexEncode(value string) string {
	intValue, err := strconv.ParseInt(value, 10, 16)
	if err != nil {
		panic(err)
	}
	stringValue := strconv.FormatInt(intValue, 16)
	return stringValue
}
