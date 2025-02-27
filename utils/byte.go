package utils

import (
	"encoding/json"
	"strconv"
)

type ByteFromString byte

func (b *ByteFromString) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	if len(s) == 1 {
		*b = ByteFromString(s[0])
		return nil
	}

	i, err := strconv.ParseUint(s, 10, 8)
	if err != nil {
		return err
	}
	*b = ByteFromString(i)
	return nil
}
