package utils

import "fmt"

func ZeroPadRight(data []byte, targetLength int) ([]byte, error) {
	dataLen := len(data)
	if dataLen > targetLength {
		return nil, fmt.Errorf("invalid length")
	}

	if dataLen == targetLength {
		return data, nil
	}

	paddedData := make([]byte, targetLength)
	copy(paddedData, data)

	return paddedData, nil
}

func ZeroPadLeft(data []byte, targetLength int) ([]byte, error) {
	dataLen := len(data)
	if dataLen > targetLength {
		return nil, fmt.Errorf("invalid length")
	}

	if dataLen == targetLength {
		return data, nil // No necesita padding
	}

	paddingNeeded := targetLength - dataLen
	paddedData := make([]byte, targetLength)

	copy(paddedData[paddingNeeded:], data)

	return paddedData, nil
}
