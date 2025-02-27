package padding

import (
	"encoding/json"
	"fmt"
)

type Type int

const (
	None Type = iota
	Fill
	Parity
)

var paddingStrings = [...]string{
	None:   "NONE",
	Fill:   "FILL",
	Parity: "PARITY",
}

// String return string
func (p *Type) String() string {
	return paddingStrings[*p]
}

// EnumIndex return index
func (p *Type) EnumIndex() int {
	return int(*p)
}

// UnmarshalJSON override default unmarshal json
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

func (p *Type) IsValid() bool {
	if int(*p) >= 0 && int(*p) < len(paddingStrings) {
		value := paddingStrings[*p]
		if value != "" {
			return true
		}
	}
	return false
}
