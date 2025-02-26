package padding

import (
	"encoding/json"
	"fmt"
)

type Position int

const (
	Right Position = iota
	Left
)

var positionStrings = [...]string{
	Right: "RIGHT",
	Left:  "LEFT",
}

// String return string
func (p *Position) String() string {
	return positionStrings[*p]
}

// EnumIndex return index
func (p *Position) EnumIndex() int {
	return int(*p)
}

// UnmarshalJSON override default unmarshal json
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

	return fmt.Errorf("invalid padding position: %s", j)
}

func (p *Position) IsValid() bool {
	if int(*p) >= 0 && int(*p) < len(positionStrings) {
		value := positionStrings[*p]
		if value != "" {
			return true
		}
	}
	return false
}
