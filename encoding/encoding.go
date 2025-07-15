// Package encoding provides various data encoding and decoding functionalities for ISO 8583 fields.
package encoding

import (
	"encoding/json"
	"fmt"
)

// Encoding represents the type of data encoding used for a field.
type Encoding int

const (
	// None indicates no specific encoding.
	None Encoding = iota
	// Bcd represents Binary-Coded Decimal encoding.
	Bcd
	// Ascii represents ASCII encoding.
	Ascii
	// Ebcdic represents EBCDIC encoding.
	Ebcdic
	// Hex represents Hexadecimal encoding.
	Hex
	// Binary represents raw binary encoding.
	Binary
)

// encodingStrings maps Encoding constants to their string representations.
var encodingStrings = [...]string{
	None:   "None",
	Bcd:    "BCD",
	Ascii:  "ASCII",
	Ebcdic: "EBCDIC",
	Hex:    "HEX",
	Binary: "BINARY",
}

// String returns the string representation of an Encoding.
func (e *Encoding) String() string {
	return encodingStrings[*e]
}

// EnumIndex returns the integer index of an Encoding.
func (e *Encoding) EnumIndex() int {
	return int(*e)
}

// UnmarshalJSON overrides the default JSON unmarshaling for Encoding.
// It allows deserializing Encoding from its string representation.
func (e *Encoding) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}

	for i, str := range encodingStrings {
		if str == j {
			*e = Encoding(i)
			return nil
		}
	}

	return fmt.Errorf("%w: %s", ErrInvalidEncodingType, j)
}

// IsValid checks if the Encoding is a valid encoding type.
func (e *Encoding) IsValid() bool {
	if int(*e) >= 0 && int(*e) < len(encodingStrings) {
		value := encodingStrings[*e]
		if value != "" {
			return true
		}
	}
	return false
}
