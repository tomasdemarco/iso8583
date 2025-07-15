package prefix

import "errors"

var (
	ErrLengthExceedsMaxLimit         = errors.New("length exceeds maximum limit")
	ErrFailedToDecodeLength          = errors.New("failed to decode length")
	ErrInvalidLengthStringConversion = errors.New("invalid length string conversion")
	ErrInvalidPrefixType             = errors.New("invalid prefix type")
)
