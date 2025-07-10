// Package padding provides functionalities for handling padding in ISO 8583 fields.
package padding

import "github.com/tomasdemarco/iso8583/encoding"

// Padder defines the interface for applying and removing padding from field values.
type Padder interface {
	// EncodePad calculates the left and right padding strings required for a given field.
	// It takes the total expected length of the field (lengthPackager), the actual data length (lengthValue),
	// and the encoder used for the field. It returns the left padding string, the right padding string,
	// and an error if padding cannot be applied (e.g., data is too long).
	EncodePad(lengthPackager int, lengthValue int, encoder encoding.Encoder) (string, string, error)
	// DecodePad calculates the number of characters to remove from the left and right
	// of a padded field value based on its total length (lengthField).
	// It returns the count of characters to remove from the left and right.
	DecodePad(lengthField int) (int, int)
	// SetChar sets the padding character for the padder.
	SetChar(char string)
	// GetChar returns the padding character used by the padder.
	GetChar() string
}

// Padders is a collection of common Padder implementations.
type Padders struct {
	NONE  Padder
	LEFT  Padder
	RIGHT Padder
}
