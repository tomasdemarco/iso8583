// Package utils provides various utility functions used across the ISO 8583 library.
package utils

import (
	"fmt"
	"strconv"
)

// Hex2Byte converts a hexadecimal string to a byte slice.
// It expects an even-length hexadecimal string.
func Hex2Byte(str string) []byte {
	slen := len(str)
	bHex := make([]byte, len(str)/2)
	ii := 0
	for i := 0; i < len(str); i = i + 2 {
		if slen != 1 {
			ss := string(str[i]) + string(str[i+1])
			bt, _ := strconv.ParseInt(ss, 16, 32)
			bHex[ii] = byte(bt)
			ii = ii + 1
			slen = slen - 2
		}
	}
	return bHex
}

// Bin2Hex converts a binary string to its hexadecimal string representation.
// It returns "error" if the input string is not a valid binary number.
func Bin2Hex(s string) string {
	ui, err := strconv.ParseUint(s, 2, 64)
	if err != nil {
		return "error"
	}
	return fmt.Sprintf("%x", ui)
}
