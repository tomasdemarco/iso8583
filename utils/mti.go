// Package utils provides various utility functions used across the ISO 8583 library.
package utils

import (
	"fmt"
	"strconv"
)

// GetMtiResponse calculates the response MTI (Message Type Indicator) for a given request MTI.
// It increments the third digit of the MTI by one to indicate a response.
// For example, "0200" becomes "0210".
// It returns the response MTI as a string.
func GetMtiResponse(mti string) (string, error) {
	if len(mti) == 4 {
		return "", ErrInvalidMTILength
	}

	v, err := strconv.Atoi(string(mti[2]))
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrInvalidMTIFormat, err)
	}

	return fmt.Sprintf("%s%d%s", mti[:2], v+1, mti[3:]), nil
}
