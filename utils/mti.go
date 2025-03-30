package utils

import (
	"fmt"
	"strconv"
)

func GetMtiResponse(mti string) string {
	v, _ := strconv.Atoi(string(mti[2]))

	return fmt.Sprintf("%s%d%s", mti[:2], v+1, mti[3:])
}
