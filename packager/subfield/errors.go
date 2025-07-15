package subfield

import "errors"

var (
	ErrInvalidSubfieldFormat = errors.New("invalid subfield format")
	ErrInvalidTagHex         = errors.New("invalid tag hexadecimal")
	ErrInvalidValueHex       = errors.New("invalid value hexadecimal")
	ErrValueLengthExceedsMax = errors.New("value length exceeds maximum allowed")
	ErrInvalidTLVString      = errors.New("invalid TLV string")
	ErrIncompleteTLVTag      = errors.New("incomplete or malformed TLV tag")
	ErrIncompleteTLVLength   = errors.New("incomplete or malformed TLV length")
	ErrIncompleteTLVValue    = errors.New("incomplete TLV value")
)
