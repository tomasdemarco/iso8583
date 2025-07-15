package prefix

import (
	"fmt"
	"math"
	"strconv"
)

func intToLenStr(length int, hex bool) string {
	if hex {
		return strconv.FormatInt(int64(length), 16)
	}

	return fmt.Sprintf("%d", length)
}

func lenStrToInt(lengthString string, hex bool) (int, error) {
	if hex {
		length, err := strconv.ParseUint(lengthString, 16, 64)
		if err != nil {
			return 0, err
		}
		return int(length), nil
	}

	length, err := strconv.Atoi(lengthString)
	if err != nil {
		return 0, err
	}
	return length, nil

}

func validateMaxLimit(length, nDigits int, hex bool) error {
	var maxLength int
	if hex {
		maxLength = int(math.Pow(16, float64(nDigits))) - 1
	} else {
		maxLength = int(math.Pow10(nDigits)) - 1
	}

	if length > maxLength {
		return fmt.Errorf("%w: length %d, max %d, digits %d", ErrLengthExceedsMaxLimit, length, maxLength, nDigits)
	}

	return nil
}
