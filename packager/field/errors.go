package field

import "errors"

var (
	ErrUnpackIndexOutOfRange = errors.New("index out of range while trying to unpack")
	ErrInvalidFieldFormat    = errors.New("invalid field format")
)
