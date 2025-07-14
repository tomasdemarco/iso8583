package packager

import "errors"

var (
	ErrFailedToOpenFile      = errors.New("failed to open packager file")
	ErrFailedToReadFile      = errors.New("failed to read packager file")
	ErrFailedToUnmarshalJSON = errors.New("failed to unmarshal packager JSON")
	ErrInvalidFieldNumber    = errors.New("invalid field number in packager JSON")
	ErrInvalidEncoding       = errors.New("invalid encoding type")
	ErrInvalidPadding        = errors.New("invalid padding type")
	ErrInvalidFieldPattern   = errors.New("invalid field pattern")
)
