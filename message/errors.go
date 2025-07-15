package message

import "errors"

var (
	ErrNotFoundInPackager       = errors.New("not found in packager")
	ErrNotFoundInMessage        = errors.New("not found in message")
	ErrMTINotFoundInPackager    = errors.New("MTI (field 0) not found in packager")
	ErrBitmapNotFoundInPackager = errors.New("bitmap (field 1) not found in packager")
)
