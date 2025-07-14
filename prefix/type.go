// Package prefix provides functionalities for handling length prefixes in ISO 8583 messages.
package prefix

import (
	"encoding/json"
	"fmt"
)

// Type represents the type of a length prefix (e.g., Fixed, L, LL, LLL).
type Type int

const (
	// Fixed indicates a fixed-length field with no explicit length prefix.
	Fixed Type = iota
	// L indicates a 1-digit length prefix.
	L
	// LL indicates a 2-digit length prefix.
	LL
	// LLL indicates a 3-digit length prefix.
	LLL
	// LLLL indicates a 4-digit length prefix.
	LLLL
	// LLLLL indicates a 5-digit length prefix.
	LLLLL
	// LLLLLL indicates a 6-digit length prefix.
	LLLLLL
)

// typeStrings maps Type constants to their string representations.
var typeStrings = [...]string{
	Fixed:  "FIXED",
	L:      "L",
	LL:     "LL",
	LLL:    "LLL",
	LLLL:   "LLLL",
	LLLLL:  "LLLLL",
	LLLLLL: "LLLLLL",
}

// String returns the string representation of a Type.
func (t *Type) String() string {
	return typeStrings[*t]
}

// EnumIndex returns the integer index of a Type.
func (t *Type) EnumIndex() int {
	return int(*t)
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

	return fmt.Errorf("%w: %s", ErrInvalidPrefixType, j)
}

// IsValid checks if the Type is a valid prefix type.
func (t *Type) IsValid() bool {
	if int(*t) >= 0 && int(*t) < len(typeStrings) {
		value := typeStrings[*t]
		if value != "" {
			return true
		}
	}
	return false
}
