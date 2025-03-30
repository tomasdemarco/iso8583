package field

import (
	"encoding/json"
	"fmt"
)

type Type int

const (
	None Type = iota
	Numeric
	String
	Binary
	Bitmap
)

var typeStrings = [...]string{
	None:    "NONE",
	Numeric: "NUMERIC",
	String:  "STRING",
	Binary:  "BINARY",
	Bitmap:  "BITMAP",
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

	return fmt.Errorf("invalid field type: %s", j)
}
