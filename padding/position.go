// Package padding provides functionalities for handling padding in ISO 8583 fields.
package padding

import (
	"encoding/json"
	"fmt"
)

// Position represents the side where padding should be applied (Right or Left).
type Position int

const (
	// Right indicates padding should be applied to the right side.
	Right Position = iota
	// Left indicates padding should be applied to the left side.
	Left
)

// positionStrings maps Position constants to their string representations.
var positionStrings = [...]string{
	Right: "RIGHT",
	Left:  "LEFT",
}

// String returns the string representation of a Position.
func (p *Position) String() string {
	return positionStrings[*p]
}

// EnumIndex returns the integer index of a Position.
func (p *Position) EnumIndex() int {
	return int(*p)
}

// UnmarshalJSON overrides the default JSON unmarshaling for Position.
// It allows deserializing Position from its string representation.
func (p *Position) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}

	for i, str := range positionStrings {
		if str == j {
			*p = Position(i)
			return nil
		}
	}

	return fmt.Errorf("%w: %s", ErrInvalidPaddingPosition, j)
}

// IsValid checks if the Position is a valid padding position.
func (p *Position) IsValid() bool {
	if int(*p) >= 0 && int(*p) < len(positionStrings) {
		value := positionStrings[*p]
		if value != "" {
			return true
		}
	}
	return false
}
