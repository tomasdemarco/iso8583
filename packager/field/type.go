// Package field defines the structure and behavior of ISO 8583 fields.
package field

import (
	"encoding/json"
	"fmt"
)

// Type represents the data type of an ISO 8583 field.
type Type int

const (
	None Type = iota
	Numeric
	String
	Binary
	Bitmap
)

// typeStrings maps Type constants to their string representations.
var typeStrings = [...]string{
	None:    "NONE",
	Numeric: "NUMERIC",
	String:  "STRING",
	Binary:  "BINARY",
	Bitmap:  "BITMAP",
}

// String returns the string representation of a Type.
func (t Type) String() string {
	return typeStrings[t]
}

// EnumIndex returns the integer index of a Type.
func (t Type) EnumIndex() int {
	return int(t)
}

// UnmarshalJSON overrides the default JSON unmarshaling for Type.
// It allows deserializing Type from its string representation.
func (t *Type) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}

	for i, str := range typeStrings {
		if str == j {
			*t = Type(i)
			return nil
		}
	}

	return fmt.Errorf("invalid field type: %s", j)
}
