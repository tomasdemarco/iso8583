package prefix

import (
	"encoding/json"
	"fmt"
)

type Type int

const (
	Fixed Type = iota
	LL
	LLL
	LLLL
)

var prefixStrings = [...]string{
	Fixed: "FIXED",
	LL:    "LL",
	LLL:   "LLL",
	LLLL:  "LLLL",
}

// String return string
func (p *Type) String() string {
	return prefixStrings[*p]
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

	for i, str := range prefixStrings {
		if str == j {
			*p = Type(i)
			return nil
		}
	}

	return fmt.Errorf("invalid prefix: %s", j)
}
