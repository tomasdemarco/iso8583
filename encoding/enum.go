package encoding

import (
	"encoding/json"
	"fmt"
)

type Encoding int

const (
	None Encoding = iota
	Bcd
	Ascii
	Ebcdic
	Hex
	Binary
)

var encodingStrings = [...]string{
	None:   "NONE",
	Bcd:    "BCD",
	Ascii:  "ASCII",
	Ebcdic: "EBCDIC",
	Hex:    "HEX",
	Binary: "BINARY",
}

// String return string
func (e *Encoding) String() string {
	return encodingStrings[*e]
}

// EnumIndex return index
func (e *Encoding) EnumIndex() int {
	return int(*e)
}

// UnmarshalJSON override default unmarshal json
func (e *Encoding) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}

	for i, str := range encodingStrings {
		if str == j {
			*e = Encoding(i)
			return nil
		}
	}

	return fmt.Errorf("invalid encoding: %s", j)
}

func (e *Encoding) IsValid() bool {
	if int(*e) >= 0 && int(*e) < len(encodingStrings) {
		value := encodingStrings[*e]
		if value != "" {
			return true
		}
	}
	return false
}
