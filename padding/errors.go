package padding

import "errors"

var (
	ErrValueTooLong           = errors.New("value too long for field")
	ErrInvalidPaddingPosition = errors.New("invalid padding position")
	ErrInvalidPaddingType     = errors.New("invalid padding type")
)
