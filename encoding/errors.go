package encoding

import "errors"

var (
	ErrNotEnoughDataToDecode = errors.New("not enough data to decode")
	ErrInvalidEncodingType   = errors.New("invalid encoding type")
)
