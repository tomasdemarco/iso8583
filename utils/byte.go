// Package utils provides various utility functions used across the ISO 8583 library.
package utils

import (
	"encoding/json"
	"strconv"
)

// ByteFromString is a custom type that allows unmarshaling a single byte
// from a JSON string, either as a character or a numeric value.
type ByteFromString byte

// UnmarshalJSON implements the json.Unmarshaler interface for ByteFromString.
// It attempts to unmarshal the JSON data as a string.
// If the string has a length of 1, it takes the first character's byte value.
// Otherwise, it attempts to parse the string as an unsigned 8-bit integer.
func (b *ByteFromString) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	if len(s) == 1 {
		*b = ByteFromString(s[0])
		return nil
	}

	i, err := strconv.ParseUint(s, 10, 8)
	if err != nil {
		return err
	}
	*b = ByteFromString(i)
	return nil
}
