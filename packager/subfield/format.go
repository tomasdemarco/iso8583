package subfield

import (
	"encoding/json"
	"fmt"
)

type Format int

const (
	None Format = iota
	Bitmap
	TLV
)

var formatStrings = [...]string{
	None:   "NONE",
	Bitmap: "BITMAP",
	TLV:    "TLV",
}

// String return string
func (f *Format) String() string {
	return formatStrings[*f]
}

// EnumIndex return index
func (f *Format) EnumIndex() int {
	return int(*f)
}

// UnmarshalJSON override default unmarshal json
func (f *Format) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}

	for i, str := range formatStrings {
		if str == j {
			*f = Format(i)
			return nil
		}
	}

	return fmt.Errorf("invalid format subfield: %s", j)
}
