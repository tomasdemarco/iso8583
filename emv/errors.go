package emv

import "errors"

var (
	ErrInvalidTagLength = errors.New("invalid tag length")
	ErrTagNotFound      = errors.New("tag not found")
)
