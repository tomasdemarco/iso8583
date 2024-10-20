package encoding

import (
	"encoding/hex"
)

func AsciiDecode(value string) (string, error) {
	asciiResult, err := hex.DecodeString(value)
	return string(asciiResult), err
}

func AsciiEncode(value string) string {
	asciiResult := hex.EncodeToString([]byte(value))
	return asciiResult
}
