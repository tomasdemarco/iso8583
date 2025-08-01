package utils

import "fmt"

// ZeroPadRight pads a byte slice with zeros on the right to reach a target length.
// It returns the padded byte slice and an error if the data length exceeds the target length.
func ZeroPadRight(data []byte, targetLength int) ([]byte, error) {
	dataLen := len(data)
	if dataLen > targetLength {
		return nil, fmt.Errorf("%w: data length %d exceeds target length %d", ErrInvalidPadLength, dataLen, targetLength)
	}

	if dataLen == targetLength {
		return data, nil // No padding needed
	}

	paddedData := make([]byte, targetLength)
	copy(paddedData, data)

	return paddedData, nil
}

// ZeroPadLeft pads a byte slice with zeros on the left to reach a target length.
// It returns the padded byte slice and an error if the data length exceeds the target length.
func ZeroPadLeft(data []byte, targetLength int) ([]byte, error) {
	dataLen := len(data)
	if dataLen > targetLength {
		return nil, fmt.Errorf("%w: data length %d exceeds target length %d", ErrInvalidPadLength, dataLen, targetLength)
	}

	if dataLen == targetLength {
		return data, nil // No padding needed
	}

	paddingNeeded := targetLength - dataLen
	paddedData := make([]byte, targetLength)

	copy(paddedData[paddingNeeded:], data)

	return paddedData, nil
}
