package prefix

import (
	"encoding/json"
	"fmt"
)

type Type int

const (
	Fixed Type = iota
	L
	LL
	LLL
	LLLL
	LLLLL
	LLLLLL
)

var typeStrings = [...]string{
	Fixed:  "FIXED",
	L:      "L",
	LL:     "LL",
	LLL:    "LLL",
	LLLL:   "LLLL",
	LLLLL:  "LLLLL",
	LLLLLL: "LLLLLL",
}

// String return string
func (t *Type) String() string {
	return typeStrings[*t]
}

// EnumIndex return index
func (t *Type) EnumIndex() int {
	return int(*t)
}

// UnmarshalJSON override default unmarshal json
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

	return fmt.Errorf("invalid type prefix: %s", j)
}
