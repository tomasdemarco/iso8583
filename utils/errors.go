package utils

import "errors"

var (
	ErrInvalidDE59Format   = errors.New("invalid DE59 format")
	ErrDE59ParsingError    = errors.New("DE59 parsing error")
	ErrInvalidMTILength    = errors.New("invalid MTI length")
	ErrInvalidMTIFormat    = errors.New("invalid MTI format")
	ErrInvalidPadLength    = errors.New("invalid padding length")
	ErrInvalidBinaryString = errors.New("invalid binary string")
)
