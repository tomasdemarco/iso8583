// Package padding provides functionalities for handling padding in ISO 8583 fields.
package padding

import (
	"encoding/json"
	"fmt"
)

// Type represents the type of padding to be applied to a field.
type Type int

const (
	// None indicates no padding should be applied.
	None Type = iota
	// Fill indicates padding with a specific character.
	Fill
	// Parity indicates padding to achieve an even/odd length.
	Parity
)

// paddingStrings maps Type constants to their string representations.
var paddingStrings = [...]string{
	None:   "NONE",
	Fill:   "FILL",
	Parity: "PARITY",
}

// String returns the string representation of a Type.
func (p *Type) String() string {
	return paddingStrings[*p]
}

// EnumIndex returns the integer index of a Type.
func (p *Type) EnumIndex() int {
	return int(*p)
}

// UnmarshalJSON overrides the default JSON unmarshaling for Type.
// It allows deserializing Type from its string representation.
func (p *Type) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}

	for i, str := range paddingStrings {
		if str == j {
			*p = Type(i)
			return nil
		}
	}

	return fmt.Errorf("invalid padding: %s", j)
}

// IsValid checks if the Type is a valid padding type.
func (p *Type) IsValid() bool {
	if int(*p) >= 0 && int(*p) < len(paddingStrings) {
		value := paddingStrings[*p]
		if value != "" {
			return true
		}
	}
	return false
}
