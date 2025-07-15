// Package prefix provides functionalities for handling length prefixes in ISO 8583 messages.
package prefix

// NonePrefixer implements the Prefixer interface for fields that do not have an explicit length prefix.
// It effectively does not encode or decode any length information.
type NonePrefixer struct {
	nDigits int
}

// NONE provides a pre-configured NonePrefixer instance for fixed-length fields.
var NONE = Prefixers{
	FIXED: &NonePrefixer{},
}

// NewNonePrefixer creates a new NonePrefixer.
// The `nDigits` parameter is typically ignored for this type of prefixer.
func NewNonePrefixer(nDigits int) Prefixer {
	return &NonePrefixer{nDigits: nDigits}
}

// EncodeLength for NonePrefixer always returns nil, nil as there is no length to encode.
func (p *NonePrefixer) EncodeLength(_ int) ([]byte, error) {
	return nil, nil
}

// DecodeLength for NonePrefixer always returns 0, nil as there is no length to decode.
func (p *NonePrefixer) DecodeLength(_ []byte, _ int) (int, error) {
	return 0, nil
}

// GetPackedLength for NonePrefixer returns the configured number of digits,
// which typically represents the fixed length of the field itself.
func (p *NonePrefixer) GetPackedLength() int {
	return p.nDigits
}
