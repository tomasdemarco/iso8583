package prefix

import (
	"fmt"
	"strconv"
)

func lengthInt(length int, hex bool) (int, error) {
	if !hex {
		return length, nil
	}

	length64, err := strconv.ParseInt(fmt.Sprintf("%d", length), 10, 16)
	if err != nil {
		return 0, err
	}

	return int(length64), nil
}

func lengthStringToInt(lengthString string, hex bool) (int, error) {
	if hex {
		length, err := strconv.ParseInt(lengthString, 16, 10)
		if err != nil {
			return 0, err
		}
		return int(length), nil
	} else {
		length, err := strconv.Atoi(lengthString)
		if err != nil {
			return 0, err
		}
		return length, nil
	}
}
