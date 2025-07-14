package bitmap

import "errors"

var (
	ErrUnpackPrimaryBitmap    = errors.New("failed to unpack primary bitmap")
	ErrInvalidPrimaryBitmap   = errors.New("invalid primary bitmap data")
	ErrUnpackSecondaryBitmap  = errors.New("failed to unpack secondary bitmap")
	ErrInvalidSecondaryBitmap = errors.New("invalid secondary bitmap data")
)
