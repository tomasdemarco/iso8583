// Package padding provides functionalities for handling padding in ISO 8583 fields.
package padding

import (
	"github.com/tomasdemarco/iso8583/encoding"
)

// NonePadder implements the Padder interface for fields that do not have any padding.
// It effectively does not apply or remove any padding.
type NonePadder struct {
	char string
}

// NONE provides a pre-configured NonePadder instance.
var NONE = Padders{
	NONE: &NonePadder{},
}

// EncodePad for NonePadder always returns empty strings and no error,
// as no padding is applied.
func (p *NonePadder) EncodePad(_ int, _ int, _ encoding.Encoder) (string, string, error) {
	return "", "", nil
}

// DecodePad for NonePadder always returns 0, 0,
// as no padding is removed.
func (p *NonePadder) DecodePad(_ int) (int, int) {
	return 0, 0
}
