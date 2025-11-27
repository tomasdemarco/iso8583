// Package field defines the structure and behavior of ISO 8583 fields.
package packager

import (
	"encoding/json"
	"fmt"
)

// FieldType represents the data type of an ISO 8583 field.
type FieldType int

const (
	None FieldType = iota
	Numeric
	String
	Binary
	Bitmap
)

// typeStrings maps FieldType constants to their string representations.
var typeStrings = [...]string{
	None:    "NONE",
	Numeric: "NUMERIC",
	String:  "STRING",
	Binary:  "BINARY",
	Bitmap:  "BITMAP",
}

// String returns the string representation of a FieldType.
func (t *FieldType) String() string {
	return typeStrings[*t]
}

// EnumIndex returns the integer index of a FieldType.
func (t *FieldType) EnumIndex() int {
	return int(*t)
}

// UnmarshalJSON overrides the default JSON unmarshaling for FieldType.
// It allows deserializing FieldType from its string representation.
func (t *FieldType) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}

	for i, str := range typeStrings {
		if str == j {
			*t = FieldType(i)
			return nil
		}
	}

	return fmt.Errorf("invalid field type: %s", j)
}
